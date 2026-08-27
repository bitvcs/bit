package domain

type Config struct {
	Url       string `json:"url"`
	BaseUrl   string `json:"base_url"`
	OrgID     string `json:"org_id"`
	ProjectID string `json:"project_id"`
}
