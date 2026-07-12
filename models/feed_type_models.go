package models



type FeedType struct {
	UID         	string `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        	string `json:"name"`
	Description 	string `json:"description"`
	Strucuture    	[]map[string]interface{} `json:"strucuture"`
}


type FeedTypeCreateRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Structure   map[string]interface{} `json:"structure"`
}




type FeedTypeUpdateRequest struct {
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	Structure   map[string]interface{} `json:"structure"`
}


