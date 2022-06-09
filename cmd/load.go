package cmd

import (
	"log"

	"github.com/fjogeleit/tracee-polr-adapter/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func loadConfig(cmd *cobra.Command) (*config.Config, error) {
	v := viper.New()

	cfgFile := ""

	configFlag := cmd.Flags().Lookup("config")
	if configFlag != nil {
		cfgFile = configFlag.Value.String()
	}

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.AddConfigPath(".")
		v.SetConfigName("config")
	}

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Println("[INFO] No configuration file found")
	}

	if flag := cmd.Flags().Lookup("kubeconfig"); flag != nil {
		v.BindPFlag("kubeconfig", flag)
	}

	if flag := cmd.Flags().Lookup("port"); flag != nil {
		v.BindPFlag("webhook.port", flag)
	}

	c := &config.Config{}

	err := v.Unmarshal(c)

	return c, err
}
