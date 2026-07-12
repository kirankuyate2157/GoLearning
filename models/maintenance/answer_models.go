package models

// QuestionAnswer represents an answer option for a question
type QuestionAnswer struct {
	UID          string `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	QuestionUID  string `json:"question_uid"`
	AnswerText   string `json:"answer_text"`
	DisplayOrder int    `json:"display_order"`
}

type QuestionAnswerCreateRequest struct {
	AnswerText   string `json:"answer_text"`
	DisplayOrder int    `json:"display_order"`
}

type QuestionAnswerUpdateRequest struct {
	AnswerText   *string `json:"answer_text"`
	DisplayOrder *int    `json:"display_order"`
}

type QuestionAnswerParams struct {
	UID         string `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	QuestionUID string `json:"question_uid"`
}
