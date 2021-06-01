// Package config handles all the user-configuration. The fields here are
// all in PascalCase but in your actual config.yml they'll be in camelCase.
// You can view the default config with `hornero --config`.
// You can open your config file by going to the status panel (using left-arrow)
// and pressing 'o'.
// You can directly edit the file (e.g. in vim) by pressing 'e' instead.
// To see the final config after your user-specific options have been merged
// with the defaults, go to the 'about' tab in the status panel.
// Because of the way we merge your user config with the defaults you may need
// to be careful: if for example you set a `commandTemplates:` yaml key but then
// give it no child values, it will scrap all of the defaults and the app will
// probably crash.
package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/OpenPeeDeeP/xdg"
	yaml "github.com/jesseduffield/yaml"
)

// UserConfig holds all of the user-configurable options
type UserConfig struct {
	// Reporting determines whether events are reported such as errors (and maybe
	// application opens but I'm not decided on that yet because it sounds kinda
	// creepy but I also would love to know how many people are using this
	// program)
	Reporting string `yaml:"reporting,omitempty"`

	// ConfirmOnQuit when enabled prompts you to confirm you want to quit when you
	// hit esc or q when no confirmation panels are open
	ConfirmOnQuit bool `yaml:"confirmOnQuit,omitempty"`

	// OS determines what defaults are set for opening files and links
	OS OSConfig `yaml:"oS,omitempty"`

	// Stats determines how long hornero will gather os stats for, and
	// what stat info to graph
	Stats StatsConfig `yaml:"stats,omitempty"`
}

// OSConfig contains config on the level of the os
type OSConfig struct {
	// OpenCommand is the command for opening a file
	OpenCommand string `yaml:"openCommand,omitempty"`

	// OpenCommand is the command for opening a link
	OpenLinkCommand string `yaml:"openLinkCommand,omitempty"`
}

// GraphConfig specifies how to make a graph of recorded os stats
type GraphConfig struct {
	// Min sets the minimum value that you want to display. If you want to set
	// this, you should also set MinType to "static". The reason for this is that
	// if Min == 0, it's not clear if it has not been set (given that the
	// zero-value of an int is 0) or if it's intentionally been set to 0.
	Min float64 `yaml:"min,omitempty"`

	// Max sets the maximum value that you want to display. If you want to set
	// this, you should also set MaxType to "static". The reason for this is that
	// if Max == 0, it's not clear if it has not been set (given that the
	// zero-value of an int is 0) or if it's intentionally been set to 0.
	Max float64 `yaml:"max,omitempty"`

	// Height sets the height of the graph in ascii characters
	Height int `yaml:"height,omitempty"`

	// Caption sets the caption of the graph. If you want to show CPU Percentage
	// you could set this to "CPU (%)"
	Caption string `yaml:"caption,omitempty"`

	// This is the path to the stat that you want to display. It is based on the
	// RecordedStats struct in os_stats.go, so feel free to look there to
	// see all the options available. Alternatively if you go into hornero and
	// go to the stats tab, you'll see that same struct in JSON format, so you can
	// just PascalCase the path and you'll have a valid path. E.g.
	// ClientStats.blkio_stats -> "ClientStats.BlkioStats"
	StatPath string `yaml:"statPath,omitempty"`

	// This determines the color of the graph. This can be any color attribute,
	// e.g. 'blue', 'green'
	Color string `yaml:"color,omitempty"`

	// MinType and MaxType are each one of "", "static". blank means the min/max
	// of the data set will be used. "static" means the min/max specified will be
	// used
	MinType string `yaml:"minType,omitempty"`

	// MaxType is just like MinType but for the max value
	MaxType string `yaml:"maxType,omitempty"`
}

// StatsConfig contains the stuff relating to stats and graphs
type StatsConfig struct {
	// Graphs contains the configuration for the stats graphs we want to show in
	// the app
	Graphs []GraphConfig

	// MaxDuration tells us how long to collect stats for. Currently this defaults
	// to "5m" i.e. 5 minutes.
	MaxDuration time.Duration `yaml:"maxDuration,omitempty"`
}

// GetDefaultConfig returns the application default configuration NOTE (to
// contributors, not users): do not default a boolean to true, because false is
// the boolean zero value and this will be ignored when parsing the user's
// config
func GetDefaultConfig() UserConfig {
	duration, err := time.ParseDuration("3m")
	if err != nil {
		panic(err)
	}

	return UserConfig{
		Reporting:     "undetermined",
		ConfirmOnQuit: false,
		OS: GetPlatformDefaultConfig(),
		Stats: StatsConfig{
			MaxDuration: duration,
			Graphs: []GraphConfig{
				{
					Caption:  "CPU (%)",
					StatPath: "DerivedStats.CPUPercentage",
					Color:    "cyan",
				},
				{
					Caption:  "Memory (%)",
					StatPath: "DerivedStats.MemoryPercentage",
					Color:    "green",
				},
			},
		},
	}
}

// AppConfig contains the base configuration fields required for hornero.
type AppConfig struct {
	Debug       bool   `long:"debug" env:"DEBUG" default:"false"`
	Version     string `long:"version" env:"VERSION" default:"unversioned"`
	Commit      string `long:"commit" env:"COMMIT"`
	BuildDate   string `long:"build-date" env:"BUILD_DATE"`
	Name        string `long:"name" env:"NAME" default:"hornero"`
	BuildSource string `long:"build-source" env:"BUILD_SOURCE" default:""`
	UserConfig  *UserConfig
	ConfigDir   string
	ProjectDir  string
}

// NewAppConfig makes a new app config
func NewAppConfig(name, version, commit, date string, buildSource string, debuggingFlag bool, composeFiles []string, projectDir string) (*AppConfig, error) {
	configDir, err := findOrCreateConfigDir(name)
	if err != nil {
		return nil, err
	}

	userConfig, err := loadUserConfigWithDefaults(configDir)
	if err != nil {
		return nil, err
	}

	// Pass compose files as individual -f flags to docker-compose
	if len(composeFiles) > 0 {
		userConfig.CommandTemplates.DockerCompose += " -f " + strings.Join(composeFiles, " -f ")
	}

	appConfig := &AppConfig{
		Name:        name,
		Version:     version,
		Commit:      commit,
		BuildDate:   date,
		Debug:       debuggingFlag || os.Getenv("DEBUG") == "TRUE",
		BuildSource: buildSource,
		UserConfig:  userConfig,
		ConfigDir:   configDir,
		ProjectDir:  projectDir,
	}

	return appConfig, nil
}

func configDirForVendor(vendor string, projectName string) string {
	envConfigDir := os.Getenv("CONFIG_DIR")
	if envConfigDir != "" {
		return envConfigDir
	}
	configDirs := xdg.New(vendor, projectName)
	return configDirs.ConfigHome()
}

func configDir(projectName string) string {
	legacyConfigDirectory := configDirForVendor("jesseduffield", projectName)
	if _, err := os.Stat(legacyConfigDirectory); !os.IsNotExist(err) {
		return legacyConfigDirectory
	}
	configDirectory := configDirForVendor("", projectName)
	return configDirectory
}

func findOrCreateConfigDir(projectName string) (string, error) {
	folder := configDir(projectName)

	err := os.MkdirAll(folder, 0755)
	if err != nil {
		return "", err
	}

	return folder, nil
}

func loadUserConfigWithDefaults(configDir string) (*UserConfig, error) {
	config := GetDefaultConfig()

	return loadUserConfig(configDir, &config)
}

func loadUserConfig(configDir string, base *UserConfig) (*UserConfig, error) {
	fileName := filepath.Join(configDir, "config.yml")

	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(fileName)
			if err != nil {
				return nil, err
			}
			file.Close()
		} else {
			return nil, err
		}
	}

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(content, base); err != nil {
		return nil, err
	}

	return base, nil
}

// WriteToUserConfig allows you to set a value on the user config to be saved
// note that if you set a zero-value, it may be ignored e.g. a false or 0 or
// empty string this is because we are using the omitempty yaml directive so
// that we don't write a heap of zero values to the user's config.yml
func (c *AppConfig) WriteToUserConfig(updateConfig func(*UserConfig) error) error {
	userConfig, err := loadUserConfig(c.ConfigDir, &UserConfig{})
	if err != nil {
		return err
	}

	if err := updateConfig(userConfig); err != nil {
		return err
	}

	file, err := os.OpenFile(c.ConfigFilename(), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	return yaml.NewEncoder(file).Encode(userConfig)
}

// ConfigFilename returns the filename of the current config file
func (c *AppConfig) ConfigFilename() string {
	return filepath.Join(c.ConfigDir, "config.yml")
}
