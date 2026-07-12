package models

type Feed struct {
	UID               string                   `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name              string                   `json:"name"`
	Code              string                   `json:"code"`
	Description       string                   `json:"description"`
	LowValue          float64                  `json:"low_value"`
	HighValue         float64                  `json:"high_value"`
	DisplayUnitUID    string                   `json:"display_unit"`
	FeedTypeUID       string                   `json:"type_id"`
	MarginOfError     float64                  `json:"margin_of_error"`
	ConnectionTypeUID string                   `json:"connection_type_uid"`
	IngestionType     string                   `json:"ingestion_type"`
	PollingRate       int                      `json:"polling_rate"`
	IsRunning         bool                     `json:"is_running"`
	JSONData          []map[string]interface{} `json:"json"`
}

// Request struct for create
type FeedCreateRequest struct {
	Name              string                   `json:"name"`
	Code              string                   `json:"code"`
	Description       string                   `json:"description"`
	LowValue          *float64                 `json:"low_value"`
	HighValue         *float64                 `json:"high_value"`
	DisplayUnitUID    string                   `json:"display_unit"`
	FeedTypeUID       string                   `json:"type_id"`
	MarginOfError     *float64                 `json:"margin_of_error"`
	ConnectionTypeUID string                   `json:"connection_type_uid"`
	IngestionType     *string                  `json:"ingestion_type"`
	PollingRate       *int                     `json:"polling_rate"`
	JSONData          []map[string]interface{} `json:"json"`
}

type FeedUpdateRequest struct {
	Name              *string                   `json:"name,omitempty"`
	Code              *string                   `json:"code,omitempty"`
	Description       *string                   `json:"description"`
	LowValue          *float64                  `json:"low_value"`
	HighValue         *float64                  `json:"high_value"`
	DisplayUnitUID    *string                   `json:"display_unit"`
	DisplayUnitID     *int                      `json:"display_unit_id"`
	FeedTypeUID       *string                   `json:"type_id"`
	FeedTypeID        *int                      `json:"feed_type_id"`
	MarginOfError     *float64                  `json:"margin_of_error"`
	ConnectionTypeUID *string                   `json:"connection_type_uid"`
	ConnectionTypeID  *int                      `json:"connection_type_id"`
	IngestionType     *string                   `json:"ingestion_type"`
	PollingRate       *int                      `json:"polling_rate"`
	JSONData          *[]map[string]interface{} `json:"json"`
	IsRunning         *bool                     `json:"is_running"`
}

type FeedMachineMappingResponse struct {
	UID        string `json:"uid"`
	FeedUID    string `json:"feed_uid"`
	MachineUID string `json:"machine_uid"`
}
