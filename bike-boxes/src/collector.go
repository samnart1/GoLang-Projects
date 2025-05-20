package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2/clientcredentials"
)

type BikeLocation struct {
	Name                    string                 `json:"name"`
	LocationID              int                    `json:"locationID"`
	Stations                []LocationStation      `json:"stations"`
	TranslatedLocationNames map[string]string      `json:"translatedLocationNames"`
	CollectedAt             time.Time              `json:"collectedAt"`
	RawData                 map[string]interface{} `json:"rawData,omitempty"`
}

type LocationStation struct {
	StationID int `json:"stationID"`
	Type      int `json:"type"`
}

type BikeStation struct {
	StationID                         int               		`json:"stationID"`
	LocationName                      string            		`json:"locationName"`
	TranslatedNames                   map[string]string 		`json:"translatedNames"`
	LocationID                        int               		`json:"locationID"`
	Name                              string            		`json:"name"`
	Address                           string            		`json:"address"`
	Addresses                         map[string]string 		`json:"addresses"`
	Latitude                          float64           		`json:"latitude"`
	Longitude                         float64           		`json:"longitude"`
	Type                              int               		`json:"type"`
	State                             int               		`json:"state"`
	CountFreePlacesAvailable_Muscular int               		`json:"countFreePlacesAvailable_MuscularBikes"`
	CountFreePlacesAvailable_Assisted int               		`json:"countFreePlacesAvailable_AssistedBikes"`
	CountFreePlacesAvailable          int               		`json:"countFreePlacesAvailable"`
	TotalPlaces                       int               		`json:"totalPlaces"`
	Places                            []Place           		`json:"places"`
	CollectedAt                       time.Time         		`json:"collectedAt"`
	RawData                           map[string]interface{}	`json:"rawData,omitempty"`
}

type Place struct {
	Position int `json:"position"`
	State    int `json:"state"`
	Level    int `json:"level"`
	Type     int `json:"type"`
}

type Collector struct {
	config        Config
	client        *http.Client
	mongoClient   *mongo.Client
	db            *mongo.Database
	locationsCol  *mongo.Collection
	stationsCol   *mongo.Collection
	scheduler     *gocron.Scheduler
	httpTransport *http.Transport
}

func NewCollector(config Config) (*Collector, error) {
	oauth2Config := &clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     fmt.Sprintf("%s/connect/token", config.APIEndpoint),
	}

	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
		IdleConnTimeout:     30 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := mongoClient.Database(config.MongoDatabase)
	locationsCol := db.Collection(config.CollectionPrefix + "_locations")
	stationsCol := db.Collection(config.CollectionPrefix + "_stations")

	locIdx := mongo.IndexModel{
		Keys: bson.D{{Key: "locationID", Value: 1}},
	}
	staIdx := mongo.IndexModel{
		Keys: bson.D{{Key: "stationID", Value: 1}},
	}
	timeIdx1 := mongo.IndexModel{
		Keys: bson.D{{Key: "collectedAt", Value: -1}},
	}
	timeIdx2 := mongo.IndexModel{
		Keys: bson.D{{Key: "collectedAt", Value: -1}},
	}

	if _, err := locationsCol.Indexes().CreateOne(ctx, locIdx); err != nil {
		log.Printf("Warning: Failed to create location index: %v", err)
	}
	if _, err := stationsCol.Indexes().CreateOne(ctx, staIdx); err != nil {
		log.Printf("Warning: Failed to create station index: %v", err)
	}
	if _, err := locationsCol.Indexes().CreateOne(ctx, timeIdx1); err != nil {
		log.Printf("Warning: Failed to create location time index: %v", err)
	}
	if _, err := stationsCol.Indexes().CreateOne(ctx, timeIdx2); err != nil {
		log.Printf("Warning: Failed to create station time index: %v", err)
	}

	collector := &Collector{
		config:        config,
		client:        oauth2Config.Client(context.Background()),
		mongoClient:   mongoClient,
		db:            db,
		locationsCol:  locationsCol,
		stationsCol:   stationsCol,
		scheduler:     gocron.NewScheduler(time.UTC),
		httpTransport: transport,
	}

	collector.client.Transport = transport

	return collector, nil
}

func (c *Collector) Start() {
	log.Println("Starting BikeBoxes Collector")

	_, err := c.scheduler.Cron(c.config.JobSchedule).Do(c.collectData)
	if err != nil {
		log.Fatalf("Failed to schedule collector job: %v", err)
	}

	c.scheduler.StartAsync()
	log.Printf("Collector job scheduled with cron: %s", c.config.JobSchedule)

	c.collectData()
}

func (c *Collector) Stop() error {
	log.Println("Stopping BikeBoxes Collector")
	c.scheduler.Stop()
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := c.mongoClient.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}
	
	log.Println("BikeBoxes Collector stopped")
	return nil
}

func (c *Collector) collectData() {
	log.Println("Starting data collection job")
	startTime := time.Now()
	
	locations, err := c.getBikeLocations()
	if err != nil {
		log.Printf("Error fetching bike locations: %v", err)
		return
	}
	
	log.Printf("Found %d bike locations", len(locations))
	
	var wg sync.WaitGroup
	stationChan := make(chan *BikeStation, 100)
	
	doneChan := make(chan bool)
	go func() {
		stationBatch := make([]*BikeStation, 0)
		batchSize := 50
		batchTimeout := time.NewTicker(5 * time.Second)
		defer batchTimeout.Stop()
		
		for {
			select {
			case station, more := <-stationChan:
				if !more {
					
					if len(stationBatch) > 0 {
						c.saveStationBatch(stationBatch)
					}
					doneChan <- true
					return
				}
				
				stationBatch = append(stationBatch, station)
				if len(stationBatch) >= batchSize {
					c.saveStationBatch(stationBatch)
					stationBatch = make([]*BikeStation, 0)
				}
				
			case <-batchTimeout.C:
				if len(stationBatch) > 0 {
					c.saveStationBatch(stationBatch)
					stationBatch = make([]*BikeStation, 0)
				}
			}
		}
	}()
	
	if err := c.saveLocations(locations); err != nil {
		log.Printf("Error saving locations to MongoDB: %v", err)
	}
	
	for _, location := range locations {
		wg.Add(1)
		go func(loc BikeLocation) {
			defer wg.Done()
			
			for _, stationRef := range loc.Stations {
				station, err := c.getBikeStation(stationRef.StationID)
				if err != nil {
					log.Printf("Error fetching station %d: %v", stationRef.StationID, err)
					continue
				}
				
				station.LocationName = loc.Name
				station.LocationID = loc.LocationID
				station.CollectedAt = time.Now()
				
				stationChan <- station
			}
		}(location)
	}
	
	wg.Wait()
	close(stationChan)
	<-doneChan
	
	duration := time.Since(startTime)
	log.Printf("Data collection job completed in %v", duration)
}

func (c *Collector) getBikeLocations() ([]BikeLocation, error) {
	
	defaultLocations, err := c.getBikeLocationsWithLanguage(c.config.DefaultLanguage)
	if err != nil {
		return nil, err
	}
	
	for i := range defaultLocations {
		defaultLocations[i].TranslatedLocationNames = make(map[string]string)
		defaultLocations[i].TranslatedLocationNames[c.config.DefaultLanguage] = defaultLocations[i].Name
		defaultLocations[i].CollectedAt = time.Now()
	}
	
	locationMap := make(map[int]*BikeLocation)
	for i := range defaultLocations {
		locationMap[defaultLocations[i].LocationID] = &defaultLocations[i]
	}
	
	for _, lang := range c.config.Languages {
		if lang == c.config.DefaultLanguage {
			continue 
		}
		
		translatedLocations, err := c.getBikeLocationsWithLanguage(lang)
		if err != nil {
			log.Printf("Warning: Failed to get locations in %s: %v", lang, err)
			continue
		}
		
		for _, translatedLoc := range translatedLocations {
			if loc, exists := locationMap[translatedLoc.LocationID]; exists {
				loc.TranslatedLocationNames[lang] = translatedLoc.Name
			}
		}
	}
	
	result := make([]BikeLocation, 0, len(locationMap))
	for _, loc := range locationMap {
		result = append(result, *loc)
	}
	
	return result, nil
}


func (c *Collector) getBikeLocationsWithLanguage(language string) ([]BikeLocation, error) {
	url := fmt.Sprintf("%s/resources/locations?languageID=%s", c.config.APIEndpoint, language)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Set("Accept", "application/json")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}
	
	var locations []BikeLocation
	if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	
	return locations, nil
}

// getBikeStation fetches a single bike station from the API
func (c *Collector) getBikeStation(stationID int) (*BikeStation, error) {
	// Get the station in the default language first
	station, err := c.getBikeStationWithLanguage(stationID, c.config.DefaultLanguage)
	if err != nil {
		return nil, err
	}
	
	// Initialize translation maps
	station.TranslatedNames = make(map[string]string)
	station.Addresses = make(map[string]string)
	
	// Add default language translations
	station.TranslatedNames[c.config.DefaultLanguage] = station.Name
	station.Addresses[c.config.DefaultLanguage] = station.Address
	
	// Get translations for all configured languages
	for _, lang := range c.config.Languages {
		if lang == c.config.DefaultLanguage {
			continue // Skip default language as we already have it
		}
		
		translatedStation, err := c.getBikeStationWithLanguage(stationID, lang)
		if err != nil {
			log.Printf("Warning: Failed to get station %d in %s: %v", stationID, lang, err)
			continue
		}
		
		// Add translations
		station.TranslatedNames[lang] = translatedStation.Name
		station.Addresses[lang] = translatedStation.Address
	}
	
	return station, nil
}

func (c *Collector) getBikeStationWithLanguage(stationID int, language string) (*BikeStation, error) {
	url := fmt.Sprintf("%s/resources/station?languageID=%s&stationID=%d", 
		c.config.APIEndpoint, language, stationID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Set("Accept", "application/json")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}
	
	var station BikeStation
	if err := json.NewDecoder(resp.Body).Decode(&station); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	
	return &station, nil
}

func (c *Collector) saveLocations(locations []BikeLocation) error {
	if len(locations) == 0 {
		return nil
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	
	docs := make([]interface{}, len(locations))
	for i, loc := range locations {
		docs[i] = loc
	}
	
	_, err := c.locationsCol.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to insert locations: %w", err)
	}
	
	log.Printf("Saved %d locations to MongoDB", len(locations))
	return nil
}

func (c *Collector) saveStationBatch(stations []*BikeStation) {
	if len(stations) == 0 {
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	docs := make([]interface{}, len(stations))
	for i, station := range stations {
		docs[i] = station
	}
	
	_, err := c.stationsCol.InsertMany(ctx, docs)
	if err != nil {
		log.Printf("Failed to insert stations: %v", err)
		return
	}
	
	log.Printf("Saved %d stations to MongoDB", len(stations))
}