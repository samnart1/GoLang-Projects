package types

type ParseResult struct {
	Data        interface{} 
	FilePath    string     
	Size        int64      
	Type        string    
	KeyCount    int        
	ArrayLength int         
}

type ValidationResult struct {
	FilePath string   
	IsValid  bool     
	Errors   []string 
	Warnings []string 
}

type OutputOptions struct {
	Format    string 
	Indent    int    
	Colors    bool   
	ShowKeys  bool   
	ShowTypes bool  
	MaxDepth  int   
	SortKeys  bool   
	UseTabs   bool 
}


type ParseRequest struct {
	Data   string `json:"data"`              
	Format string `json:"format,omitempty"` 
	Indent int    `json:"indent,omitempty"`  
}


type ParseResponse struct {
	Success   bool        `json:"success"`       
	Data      interface{} `json:"data,omitempty"`  
	Formatted string      `json:"formatted,omitempty"` 
	Error     string      `json:"error,omitempty"`  
	Metadata  *Metadata   `json:"metadata,omitempty"` 
}


type ValidationResponse struct {
	Success  bool     `json:"success"`          
	Valid    bool     `json:"valid"`            
	Errors   []string `json:"errors,omitempty"` 
	Warnings []string `json:"warnings,omitempty"` 
	Error    string   `json:"error,omitempty"`  
}


type Metadata struct {
	Type        string `json:"type"`         
	Size        int64  `json:"size"`         
	KeyCount    int    `json:"keyCount,omitempty"` 
	ArrayLength int    `json:"arrayLength,omitempty"`
}