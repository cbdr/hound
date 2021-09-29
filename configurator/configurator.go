package configurator

import (
	"encoding/json"
	"os"
)

type Repo struct {
	Url             string `json:"url"`
	MsBetweenPolls  int    `json:"ms-between-poll"`
	ExcludeDotFiles bool   `json:"exclude-dot-files"`
}

type Config struct {
	DbPath                string           `json:"dbpath"`
	Title                 string           `json:"title"`
	Repos                 map[string]*Repo `json:"repos"`
	MaxConcurrentIndexers int              `json:"max-concurrent-indexers"`
	HealthCheckURI        string           `json:"health-check-uri"`
	VcsConfg              map[string]*Vcs  `json:"vcs-config"`
}

type Vcs struct {
	DetectRef bool `json:"detect-ref"`
}

func (c *Config) SaveToFile(filename string) error {
	w, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer w.Close()

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(c); err != nil {
		return err
	}

	return nil
}
