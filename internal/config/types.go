package config

// TargetMode defines how the persona is applied to the target
type TargetMode string

const (
	TargetModeLink TargetMode = "link"
	TargetModeCopy TargetMode = "copy"
)

// TargetConfig represents a single target in the configuration
type TargetConfig struct {
	Path string     `mapstructure:"path"`
	Mode TargetMode `mapstructure:"mode"`
}

// Config represents the top-level configuration
type Config struct {
	AgentsDir  []string       `mapstructure:"agents_dir" yaml:"agents_dir"`
	TargetFile string         `mapstructure:"target_file" yaml:"target_file"` // Legacy support
	Targets    []TargetConfig `mapstructure:"targets" yaml:"targets"`
}
