package models

type AlertsTypesGet struct {
	UID         string                   `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Structure   []map[string]interface{} `json:"structure"`
}

type AlertsGet struct {
	UID     string                   `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name    string                   `json:"name"`
	TypeUID string                   `json:"typeUID"  example:"550e8400-e29b-41d4-a716-446655440000"`
	Config  []map[string]interface{} `json:"structure"`
}

type AlertsPost struct {
	Name      string                   `json:"name"`
	TypeUID   string                   `json:"typeUID"  example:"550e8400-e29b-41d4-a716-446655440000"`
	Structure []map[string]interface{} `json:"structure"`
}

type AlertsPatch struct {
	Name      *string                   `json:"name"`
	TypeUID   *string                   `json:"typeUID"  example:"550e8400-e29b-41d4-a716-446655440000"`
	Structure *[]map[string]interface{} `json:"structure"`
}

type AlertsTypeParams struct {
	UID *string `json:"typeUID"  example:"550e8400-e29b-41d4-a716-446655440000"`
}

type AlertsParams struct {
	UID *string `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
}
