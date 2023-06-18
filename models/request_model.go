package models

type Request struct {
	Index   string  `json:"index"`
	Records []Email `json:"records"`
}
