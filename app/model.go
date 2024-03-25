package main

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ParticipantDetailRequest struct {
	Username   string `json:"username"`
	ExternalId string `json:"externalid"`
}

type SubmitDeferralRequest struct {
	Username        string            `json:"username"`
	ExternalPlanId  string            `json:"planid"`
	DeductMethod    int               `json:"deductmethod"`
	DeferralRequest []DeferralRequest `json:"deferral"`
}

type Source struct {
	SourceName string `json:"source"`
	SourceType string `json:"sourcetype"`
}

type PlanRequest struct {
	ExternalPlanId string `json:"planid"`
}

type Plan struct {
	ExternalPlanId string `json:"planid"`
	PlanName       string `json:"planname"`
}

type DeferralRequest struct {
	Source         string `json:"source"`
	ExternalPlanId string `json:"planid"`
	DeductAmount   int    `json:"deductamount"`
}

type Customer struct {
	Name       string
	Customerid string
	Password   string
}

type Deferral struct {
	Source       Source `json:"source"`
	Plan         Plan   `json:"planid"`
	DeductAmount int    `json:"deductamount"`
	DeductMethod int    `json:"deductmethod"`
}
type ParticipantPlanDetailResponse struct {
	ExternalPlanId string
	PlanName       string
	Deferrals      map[string]*Deferral //sourcename -> Deferral
}
