package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	ConcurrentUsers int
	TestDuration    string
	RequestRate     int
	TargetURL       string
	HTTPMethod      string
	CustomHeaders   map[string]string
	RequestPayload  string
	Timeout         int
	KeepAlive       bool
	TLSSkipVerify   bool
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) AddFlags(cmd *cobra.Command) {
	cmd.Flags().IntVarP(&c.ConcurrentUsers, "users", "u", 1, "Number of concurrent users")
	cmd.Flags().StringVarP(&c.TestDuration, "duration", "d", "1m", "Test duration (e.g., 1s, 1m, 1h)")
	cmd.Flags().IntVarP(&c.RequestRate, "rate", "r", 1, "Request rate (requests per second)")
	cmd.Flags().StringVarP(&c.TargetURL, "url", "l", "", "Target URL")
	cmd.Flags().StringVarP(&c.HTTPMethod, "method", "m", "GET", "HTTP method")
	cmd.Flags().StringToStringP("headers", "H", nil, "Custom headers (e.g., -H 'Content-Type=application/json')")
	cmd.Flags().StringVarP(&c.RequestPayload, "payload", "p", "", "Request payload for POST/PUT requests")
	cmd.Flags().IntVar(&c.Timeout, "timeout", 30, "Request timeout in seconds")
	cmd.Flags().BoolVar(&c.KeepAlive, "keepalive", true, "Use HTTP keep-alive")
	cmd.Flags().BoolVar(&c.TLSSkipVerify, "insecure", false, "Skip TLS certificate verification")

	cmd.MarkFlagRequired("url")
}

func (c *Config) LoadConfig() error {
	c.CustomHeaders = viper.GetStringMapString("headers")
	return nil
}
