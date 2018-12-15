package prometheus

// ScrapeConfig is a minimalist version of the Prometheus ScrapeConfig structure
type ScrapeConfig struct {
	JobName        string            `json:"job_name"`
	HonorLabels    bool              `json:"honor_labels,omitempty"`
	ScrapeInterval string            `json:"scrape_interval,omitempty"`
	ScrapeTimeout  string            `json:"scrape_timeout,omitempty"`
	MetricsPath    string            `json:"metrics_path,omitempty"`
	Scheme         string            `json:"scheme,omitempty"`
	StaticConfigs  []StaticConfig    `json:"static_configs,omitempty"`
	BasicAuth      map[string]string `json:"basic_auth,omitempty"`
}

// StaticConfig is a minimalist version of the Prometheus StaticConfig structure
type StaticConfig struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels,omitempty"`
}
