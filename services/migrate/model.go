package migrate

type Migrations struct {
	ForceVersion int    `json:"forceVersion,omitempty"`
	Type         string `json:"type,omitempty"`
}

type Request struct {
	Migrations []Migrations `json:"migrations,omitempty"`
}
