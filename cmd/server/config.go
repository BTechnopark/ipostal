package main

type Config struct {
	BaseUrl string
}

// GetBaseUrl implements api.PostalConfig.
func (c *Config) GetBaseUrl() string {
	return c.BaseUrl
}
