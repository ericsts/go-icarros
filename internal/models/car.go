package models

type Car struct {
	ID     int     `json:"id"`
	UserID int     `json:"user_id"`
	Marca  string  `json:"marca"`
	Modelo string  `json:"modelo"`
	Ano    int     `json:"ano"`
	Valor  float64 `json:"valor"`
}
