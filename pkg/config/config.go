package config

// Profiling configuration
type Profiling struct {
	Enabled bool `mapstructure:"enabled"`
}

// Webhook configuration
type Webhook struct {
	Port int `mapstructure:"port"`
}

// Results configuration
type Results struct {
	MaxPerReport    int `mapstructure:"maxPerReport"`
	MinimumSeverity int `mapstructure:"minimumSeverity"`
}

// Results configuration
type Rules struct {
	Exclude []string `mapstructure:"exclude"`
}

// Config of the Tracee Adapter
type Config struct {
	Kubeconfig string    `mapstructure:"kubeconfig"`
	Profiling  Profiling `mapstructure:"profiling"`
	Webhook    Webhook   `mapstructure:"webhook"`
	Results    Results   `mapstructure:"results"`
	Rules      Rules     `mapstructure:"rules"`
}
