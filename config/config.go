package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type HTTPConfig struct {
	Address      string        `yaml:"address"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	IdleTimeout  time.Duration `yaml:"idleTimeout"`
}

type Config struct {
	HTTP     HTTPConfig     `yaml:"http"`
	Upstream UpstreamConfig `yaml:"upstream"`
	Routes   []Route        `yaml:"routes"`
	Pools    []Pool         `yaml:"pools"`
	CORS     CORSConfig     `yaml:"cors"`
	Auth     AuthConfig     `yaml:"auth"`
	Rate     RateConfig     `yaml:"rate"`
}

type UpstreamConfig struct {
	URL     string        `yaml:"url"`
	Timeout time.Duration `yaml:"timeout"`
}

type Pool struct {
	Name     string           `yaml:"name"`
	Strategy string           `yaml:"strategy"` // rr | least_conn
	Backends []UpstreamConfig `yaml:"backends"`
	Health   HealthCheck      `yaml:"health"`
}

type HealthCheck struct {
	Path            string        `yaml:"path"`
	Interval        time.Duration `yaml:"interval"`
	Timeout         time.Duration `yaml:"timeout"`
	UnhealthyThresh int           `yaml:"unhealthyThreshold"`
	HealthyThresh   int           `yaml:"healthyThreshold"`
}

type Route struct {
	Name        string            `yaml:"name"`
	PathPrefix  string            `yaml:"pathPrefix"`
	Methods     []string          `yaml:"methods"`
	Headers     map[string]string `yaml:"headers"`
	Upstream    UpstreamConfig    `yaml:"upstream"`
	Pool        string            `yaml:"pool"`
	StripPrefix bool              `yaml:"stripPrefix"`
}

type CORSConfig struct {
	Enabled          bool     `yaml:"enabled"`
	AllowedOrigins   []string `yaml:"allowedOrigins"`
	AllowedMethods   []string `yaml:"allowedMethods"`
	AllowedHeaders   []string `yaml:"allowedHeaders"`
	ExposeHeaders    []string `yaml:"exposeHeaders"`
	AllowCredentials bool     `yaml:"allowCredentials"`
	MaxAge           int      `yaml:"maxAge"`
}

type AuthConfig struct {
	APIKeys []string  `yaml:"apiKeys"`
	JWT     JWTConfig `yaml:"jwt"`
}

type JWTConfig struct {
	Issuer   string   `yaml:"issuer"`
	Audience []string `yaml:"audience"`
	Secret   string   `yaml:"secret"`
	Enabled  bool     `yaml:"enabled"`
}

type RateConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Requests int           `yaml:"requests"`
	Per      time.Duration `yaml:"per"`
}

func defaultConfig() Config {
	return Config{
		HTTP: HTTPConfig{
			Address:      ":8081",
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Upstream: UpstreamConfig{
			URL:     "http://localhost:8080",
			Timeout: 15 * time.Second,
		},
		Routes: []Route{},
		Pools:  []Pool{},
		CORS: CORSConfig{
			Enabled:          true,
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			ExposeHeaders:    []string{},
			AllowCredentials: false,
			MaxAge:           600,
		},
		Auth: AuthConfig{
			APIKeys: []string{},
			JWT:     JWTConfig{Enabled: false, Secret: ""},
		},
		Rate: RateConfig{Enabled: false, Requests: 100, Per: time.Minute},
	}
}

func Load() Config {
	cfg := defaultConfig()
	path := os.Getenv("GATEWAY_CONFIG")
	if path == "" {
		return cfg
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}
	_ = yaml.Unmarshal(b, &cfg)
	return cfg
}
