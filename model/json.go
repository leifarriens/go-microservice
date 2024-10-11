package model

type PageableRequest struct {
	Limit  int `query:"limit" validate:"number,min=1,max=100"`
	Offset int `query:"offset" validate:"number,min=0,max=100"`
}

type PageableResponse struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
