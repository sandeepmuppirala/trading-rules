package models

// DataRulesInput : Input dto
type DataRulesInput struct {
	ID             string `json:"id"`
	CustomerID     string `json:"customer_id"`
	LoadAmount     string `json:"load_amount"`
	Time           string `json:"time"`
	AttemptsPerDay int
}
