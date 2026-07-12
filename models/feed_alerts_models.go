package models

type FeedAlertGet struct {
	UID            string  `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	AlertUID       string  `json:"alert_uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	LowAlertValue  float64 `json:"low_alert_value"`
	HighAlertValue float64 `json:"high_alert_value"`
	FeedUID        string  `json:"feed_uid" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type FeedAlertPost struct {
	AlertUID       string  `json:"alert_uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	LowAlertValue  float64 `json:"low_alert_value" validate:"required"`
	HighAlertValue float64 `json:"high_alert_value" validate:"required"`
}

type FeedAlertPatch struct {
	AlertUID       *string  `json:"alert_uid,omitempty"  example:"550e8400-e29b-41d4-a716-446655440000"`
	LowAlertValue  *float64 `json:"low_alert_value,omitempty"`
	HighAlertValue *float64 `json:"high_alert_value,omitempty"`
}
