package models

type ResponseZinc struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source struct {
				Email
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
