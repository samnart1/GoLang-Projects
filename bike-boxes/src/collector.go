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

// BikeLocation represents a location with bike stations
type BikeLocation struct {
	Name                    string                 `json:"name"`
	LocationID              int                    `json:"locationID"`
	Stations                []LocationStation      `json:"stations"`
	TranslatedLocationNames map[string]string      `json:"translatedLocationNames"`
	CollectedAt             time.Time              `json:"collectedAt"`
	RawData                 map[string]interface{} `json:"rawData,omitempty"`
}

// LocationStation represents a station within a location
type LocationStation struct {
	StationID int `json:"stationID"`
	Type      int `json:"type"`
}

// BikeStation represents a bike station with availability information
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

// Place represents a bike place/slot within a station
type Place struct {
	Position int `json:"position"`
	State    int `json:"state"`
	Level    int `json:"level"`
	Type     int `json:"type"`
}

// Collector handles data collection from bike data API
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

// NewCollector creates a new collector instance
func NewCollector(config Config) (*Collector, error) {
	oauth2Config := &clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     fmt.Sprintf("%s/connect/token", config.APIEndpoint),
	}

	// Create custom transport with reasonable timeouts
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
		IdleConnTimeout:     30 * time.Second,
	}

	// Setup context for MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping MongoDB to verify connection
	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Setup collections
	db := mongoClient.Database(config.MongoDatabase)
	locationsCol := db.Collection(config.CollectionPrefix + "_locations")
	stationsCol := db.Collection(config.CollectionPrefix + "_stations")

	// Create indexes
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

	// Pass transport to the OAuth client
	collector.client.Transport = transport

	return collector, nil
}

// Start begins the collector job scheduler
func (c *Collector) Start() {
	log.Println("Starting BikeBoxes Collector")

	// Schedule the collector job
	_, err := c.scheduler.Cron(c.config.JobSchedule).Do(c.collectData)
	if err != nil {
		log.Fatalf("Failed to schedule collector job: %v", err)
	}

	// Start the scheduler
	c.scheduler.StartAsync()
	log.Printf("Collector job scheduled with cron: %s", c.config.JobSchedule)

	// Run once at startup
	c.collectData()
}

// Stop gracefully stops the collector
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

// collectData is the main function that collects all data
func (c *Collector) collectData() {
	log.Println("Starting data collection job")
	startTime := time.Now()
	
	// Get the bike locations first
	locations, err := c.getBikeLocations()
	if err != nil {
		log.Printf("Error fetching bike locations: %v", err)
		return
	}
	
	log.Printf("Found %d bike locations", len(locations))
	
	// Process each location and get its stations
	var wg sync.WaitGroup
	stationChan := make(chan *BikeStation, 100)
	
	// Start MongoDB writer goroutine
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
					// Channel closed, process any remaining stations
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
	
	// Save the locations to MongoDB
	if err := c.saveLocations(locations); err != nil {
		log.Printf("Error saving locations to MongoDB: %v", err)
	}
	
	// Process each location to get its stations
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
				
				// Set location name
				station.LocationName = loc.Name
				station.LocationID = loc.LocationID
				station.CollectedAt = time.Now()
				
				// Send the station to the channel
				stationChan <- station
			}
		}(location)
	}
	
	// Wait for all goroutines to finish
	wg.Wait()
	close(stationChan)
	<-doneChan
	
	duration := time.Since(startTime)
	log.Printf("Data collection job completed in %v", duration)
}

// getBikeLocations fetches all bike locations from the API
func (c *Collector) getBikeLocations() ([]BikeLocation, error) {
	// Get locations using the default language first
	defaultLocations, err := c.getBikeLocationsWithLanguage(c.config.DefaultLanguage)
	if err != nil {
		return nil, err
	}
	
	// Initialize translatedLocationNames for each location
	for i := range defaultLocations {
		defaultLocations[i].TranslatedLocationNames = make(map[string]string)
		defaultLocations[i].TranslatedLocationNames[c.config.DefaultLanguage] = defaultLocations[i].Name
		defaultLocations[i].CollectedAt = time.Now()
	}
	
	// Map locations by ID for easy access
	locationMap := make(map[int]*BikeLocation)
	for i := range defaultLocations {
		locationMap[defaultLocations[i].LocationID] = &defaultLocations[i]
	}
	
	// Get translations for all configured languages
	for _, lang := range c.config.Languages {
		if lang == c.config.DefaultLanguage {
			continue // Skip default language as we already have it
		}
		
		translatedLocations, err := c.getBikeLocationsWithLanguage(lang)
		if err != nil {
			log.Printf("Warning: Failed to get locations in %s: %v", lang, err)
			continue
		}
		
		// Add translations to the map
		for _, translatedLoc := range translatedLocations {
			if loc, exists := locationMap[translatedLoc.LocationID]; exists {
				loc.TranslatedLocationNames[lang] = translatedLoc.Name
			}
		}
	}
	
	// Convert map back to slice
	result := make([]BikeLocation, 0, len(locationMap))
	for _, loc := range locationMap {
		result = append(result, *loc)
	}
	
	return result, nil
}

// getBikeLocationsWithLanguage fetches bike locations for a specific language
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

// getBikeStationWithLanguage fetches a station with a specific language
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

// saveLocations saves a batch of locations to MongoDB
func (c *Collector) saveLocations(locations []BikeLocation) error {
	if len(locations) == 0 {
		return nil
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Convert to interface{} slice for MongoDB
	docs := make([]interface{}, len(locations))
	for i, loc := range locations {
		docs[i] = loc
	}
	
	// Insert into MongoDB
	_, err := c.locationsCol.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to insert locations: %w", err)
	}
	
	log.Printf("Saved %d locations to MongoDB", len(locations))
	return nil
}

// saveStationBatch saves a batch of stations to MongoDB
func (c *Collector) saveStationBatch(stations []*BikeStation) {
	if len(stations) == 0 {
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Convert to interface{} slice for MongoDB
	docs := make([]interface{}, len(stations))
	for i, station := range stations {
		docs[i] = station
	}
	
	// Insert into MongoDB
	_, err := c.stationsCol.InsertMany(ctx, docs)
	if err != nil {
		log.Printf("Failed to insert stations: %v", err)
		return
	}
	
	log.Printf("Saved %d stations to MongoDB", len(stations))
}