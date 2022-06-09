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

// Config of the Tracee Adapter
type Config struct {
	Kubeconfig string    `mapstructure:"kubeconfig"`
	Profiling  Profiling `mapstructure:"profiling"`
	Webhook    Webhook   `mapstructure:"webhook"`
	Results    Results   `mapstructure:"results"`
}
