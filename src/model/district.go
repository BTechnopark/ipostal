package model

type District struct {
	ID       int    `json:"id"`
	Province string `json:"province"`
	Regency  string `json:"regency"`
	District string `json:"district"`
}
