package prometheus

// ScrapeConfig is a minimalist version of the Prometheus ScrapeConfig structure
type ScrapeConfig struct {
	JobName         string            `yaml:"job_name"`
	HonorLabels     bool              `yaml:"honor_labels,omitempty"`
	Params          interface{}       `yaml:"params,omitempty"`
	ScrapeInterval  string            `yaml:"scrape_interval,omitempty"`
	ScrapeTimeout   string            `yaml:"scrape_timeout,omitempty"`
	MetricsPath     string            `yaml:"metrics_path,omitempty"`
	Scheme          string            `yaml:"scheme,omitempty"`
	StaticConfigs   []StaticConfig    `yaml:"static_configs,omitempty"`
	BasicAuth       map[string]string `yaml:"basic_auth,omitempty"`
	BearerToken     string            `yaml:"bearer_token,omitempty"`
	BearerTokenFile string            `yaml:"bearer_token_file,omitempty"`
	ProxyURL        string            `yaml:"proxy_url,omitempty"`
	TLSConfig       interface{}       `yaml:"tls_config,omitempty"`

	DNSSDConfigs        interface{} `yaml:"dns_sd_configs,omitempty"`
	FileSDConfigs       interface{} `yaml:"file_sd_configs,omitempty"`
	ConsulSDConfigs     interface{} `yaml:"consul_sd_configs,omitempty"`
	ServersetSDConfigs  interface{} `yaml:"serverset_sd_configs,omitempty"`
	NerveSDConfigs      interface{} `yaml:"nerve_sd_configs,omitempty"`
	MarathonSDConfigs   interface{} `yaml:"marathon_sd_configs,omitempty"`
	KubernetesSDConfigs interface{} `yaml:"kubernetes_sd_configs,omitempty"`
	GCESDConfigs        interface{} `yaml:"gce_sd_configs,omitempty"`
	EC2SDConfigs        interface{} `yaml:"ec2_sd_configs,omitempty"`
	OpenstackSDConfigs  interface{} `yaml:"openstack_sd_configs,omitempty"`
	AzureSDConfigs      interface{} `yaml:"azure_sd_configs,omitempty"`
	TritonSDConfigs     interface{} `yaml:"triton_sd_configs,omitempty"`

	RelabelConfigs     interface{} `yaml:"relabel_configs,omitempty"`
	MetricLabelConfigs interface{} `yaml:"metric_relabel_configs,omitempty"`
}

// StaticConfig is a minimalist version of the Prometheus StaticConfig structure
type StaticConfig struct {
	Targets []string          `yaml:"targets"`
	Labels  map[string]string `yaml:"labels,omitempty"`
}
