package prometheus

// ScrapeConfig is a minimalist version of the Prometheus ScrapeConfig structure
type ScrapeConfig struct {
	JobName         string            `json:"job_name"`
	HonorLabels     bool              `json:"honor_labels,omitempty"`
	Params          interface{}       `json:"params,omitempty"`
	ScrapeInterval  string            `json:"scrape_interval,omitempty"`
	ScrapeTimeout   string            `json:"scrape_timeout,omitempty"`
	MetricsPath     string            `json:"metrics_path,omitempty"`
	Scheme          string            `json:"scheme,omitempty"`
	StaticConfigs   []StaticConfig    `json:"static_configs,omitempty"`
	BasicAuth       map[string]string `json:"basic_auth,omitempty"`
	BearerToken     string            `json:"bearer_token,omitempty"`
	BearerTokenFile string            `json:"bearer_token_file,omitempty"`
	ProxyURL        string            `json:"proxy_url,omitempty"`
	TLSConfig       interface{}       `json:"tls_config,omitempty"`

	DNSSDConfigs        interface{} `json:"dns_sd_configs,omitempty"`
	FileSDConfigs       interface{} `json:"file_sd_configs,omitempty"`
	ConsulSDConfigs     interface{} `json:"consul_sd_configs,omitempty"`
	ServersetSDConfigs  interface{} `json:"serverset_sd_configs,omitempty"`
	NerveSDConfigs      interface{} `json:"nerve_sd_configs,omitempty"`
	MarathonSDConfigs   interface{} `json:"marathon_sd_configs,omitempty"`
	KubernetesSDConfigs interface{} `json:"kubernetes_sd_configs,omitempty"`
	GCESDConfigs        interface{} `json:"gce_sd_configs,omitempty"`
	EC2SDConfigs        interface{} `json:"ec2_sd_configs,omitempty"`
	OpenstackSDConfigs  interface{} `json:"openstack_sd_configs,omitempty"`
	AzureSDConfigs      interface{} `json:"azure_sd_configs,omitempty"`
	TritonSDConfigs     interface{} `json:"triton_sd_configs,omitempty"`

	RelabelConfigs     interface{} `json:"relabel_configs,omitempty"`
	MetricLabelConfigs interface{} `json:"metric_relabel_configs,omitempty"`
}

// StaticConfig is a minimalist version of the Prometheus StaticConfig structure
type StaticConfig struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels,omitempty"`
}
