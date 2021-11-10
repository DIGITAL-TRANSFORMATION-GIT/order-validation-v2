package models

type Requirements struct {
	Id              int    `json:"id,omitempty"`
	Request         string `json:"request"`
	ExpectedOutcome string `json:"outcome"`
	Status          bool   `json:"status"`
	OrderID         string `json:"order_id,omitempty"`
	UserID          string `json:"user_id"`
}

type RequirementPatch struct {
	Patches []Patch `json:"patch"`
}

type Patch struct {
	Id              int     `json:"id"`
	UserID          *string `json:"user_id,omitempty"`
	ExpectedOutcome *string `json:"outcome,omitempty"`
}