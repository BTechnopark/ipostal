package model

type PostalCode struct {
	PostalCode string `json:"postal_code"`
	Province   string `json:"province"`
	Region     string `json:"region"`
	District   string `json:"district"`
	Village    string `json:"village"`
}
type ListPostalCode []*PostalCode
