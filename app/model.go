package main

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ParticipantDetailRequest struct {
	Username   string `json:"username"`
	ExternalId string `json:"externalid"`
}

type Source struct {
	SourceName string
	SourceType string
}

type Plan struct {
	ExternalPlanId string
	PlanName       string
}

type Deferral struct {
	Source       Source
	Plan         Plan
	DeductAmount int
	DeductMethod int
}

type ParticipantPlanDetailResponse struct {
	ExternalPlanId string
	PlanName       string
	Deferrals      map[string]*Deferral //sourcename -> Deferral
}
