package model

type PostalCode struct {
	ID         int    `json:"id"`
	PostalCode string `json:"postal_code"`
	Province   string `json:"province"`
	Regency    string `json:"regency"`
	District   string `json:"district"`
	Village    string `json:"village"`
}
