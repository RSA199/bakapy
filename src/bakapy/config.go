package bakapy

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	IncludeJobs []string `yaml:"include_jobs"`
	Listen      string
	StorageDir  string    `yaml:"storage_dir"`
	MetadataDir string    `yaml:"metadata_dir"`
	StatusDir   string    `yaml:"status_dir"`
	Ports       PortRange `yaml:"port_range"`
	Options     GlobalOptions
	Jobs        map[string]JobConfig
}

func (cfg *Config) PrettyFmt() []byte {
	s, _ := yaml.Marshal(cfg)
	return s
}

type GlobalOptions struct {
	Gzip bool
	Args map[string]string
}

type PortRange struct {
	Start int
	End   int
}

func (p *PortRange) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value map[string]int
	err := unmarshal(&value)
	if err != nil {
		panic(err)
	}

	if value["start"] >= value["end"] {
		return errors.New("port_range.start cannot be more or equal than port_range.end")
	}
	p.Start = value["start"]
	p.End = value["end"]
	return nil
}

type RunAtSpec struct {
	Minute  string
	Hour    string
	Day     string
	Month   string
	Weekday string
}

func (r *RunAtSpec) SchedulerString() string {
	return fmt.Sprintf(
		"20 %s %s %s %s %s",
		r.Minute,
		r.Hour,
		r.Day,
		r.Month,
		r.Weekday,
	)
}

type JobConfig struct {
	Disabled   bool
	MaxAgeDays int `yaml:"max_age_days"`
	Namespace  string
	Host       string
	Command    string
	Args       map[string]string
	RunAt      RunAtSpec `yaml:"run_at"`
}

func NewConfig() *Config {
	jobs := Config{
		Jobs: map[string]JobConfig{},
	}
	return &jobs
}

func ParseConfig(configPath string) (*Config, error) {
	cfg := NewConfig()

	rawConfig, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(rawConfig, cfg)
	if err != nil {
		return nil, err
	}

	var jobs = map[string]JobConfig{}
	for _, relPath := range cfg.IncludeJobs {

		raw, err := ioutil.ReadFile(relPath)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(raw, &jobs)
		if err != nil {
			return nil, err
		}
		for name, params := range jobs {
			cfg.Jobs[name] = params
		}

	}

	return cfg, err
}