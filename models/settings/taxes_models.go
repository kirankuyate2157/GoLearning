package models

type TaxesGet struct {
	UID        string  `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	TYPE       string  `json:"type"`
	PERCENTAGE float64 `json:"percentage"`
}

type TaxesPost struct {
	Type       string  `json:"type"`
	Percentage float64 `json:"percentage"`
}

type TaxesPatch struct {
	Type       *string  `json:"type"`
	Percentage *float64 `json:"percentage"`
}
