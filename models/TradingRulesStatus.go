package models

// DataRulesOutput : Output DTO
type DataRulesOutput struct {
	ID          string `json:"id"`
	CustomerID 	string `json:"customer_id"`
	Accepted    bool   `json:"accepted"`
}
