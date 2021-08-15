package types

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"path"
)

type ConfigSite struct {
	Name          string `yaml:"name" json:"name"`
	Enabled       bool   `yaml:"enabled" json:"enabled"`
	LabelSelector string `yaml:"labelSelector" json:"labelSelector"`
	Namespace     string `yaml:"namespace" json:"namespace"`
	User          struct {
		Email    string `yaml:"email" json:"name"`
		Username string `yaml:"username" json:"username"`
		Password string `yaml:"password" json:"password"`
	} `yaml:"user" json:"user"`
}

type DataUrl struct {
	Name string `yaml:"name" json:"name"`
	Path string `yaml:"path" json:"path"`
}

type Wordpress struct {
	StoreRoot string       `yaml:"storeRoot" envconfig:"WP_STORE_ROOT" json:"storeRoot"`
	Sites     []ConfigSite `yaml:"sites" json:"sites"`
	DataUrls  []DataUrl    `yaml:"dataUrls" json:"dataUrls"`
	Themes    []Theme      `yaml:"themes" json:"themes"`
	Plugins   []Plugin     `yaml:"plugins" json:"plugins"`
}

type Config struct {
	K8RestConfig *rest.Config `json:"-"`
	// Runtime not meant for env or yaml files, only flags
	Runtime struct {
		SpecificPod string
	} `json:"-"`

	Environment string `yaml:"environment" envconfig:"ENVIRONMENT" json:"environment"`
	DatabaseDsn string `yaml:"databaseDsn" envconfig:"DATABASE_DSN" json:"-"`

	K8 struct {
		LabelSelector string `yaml:"labelSelector" envconfig:"K8_LABEL_SELECTOR" json:"labelSelector"`
	} `yaml:"kubernetes" json:"kubernetes"`

	Logging struct {
		Level string `yaml:"level" envconfig:"LOGGING_LEVEL" json:"level"`
	} `yaml:"logging" json:"logging"`

	ManifestUrl string    `yaml:"manifest" envconfig:"WP_MANIFEST" json:"manifest"`
	Wordpress   Wordpress `yaml:"wordpress" json:"wordpress"`
}

func (config *Config) GetSite(name string) *ConfigSite {
	for _, site := range config.Wordpress.Sites {
		if site.Name == name {
			return &site
		}
	}

	log.Fatalf("Could not find site name: %s", name)

	return nil
}

func (wp *Wordpress) GetPluginByName(name string) (*Plugin, error) {
	for _, plugin := range wp.Plugins {
		if plugin.Name == name {
			return &plugin, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("could not find plugin %s", name))
}

func (wp *Wordpress) HasPluginByName(name string) bool {
	_, err := wp.GetPluginByName(name)

	return err == nil
}

func (wp *Wordpress) AddPlugin(plugin Plugin) {
	wp.Plugins = append(wp.Plugins, plugin)
}

func (wp *Wordpress) GetThemeByName(name string) (*Theme, error) {
	for _, theme := range wp.Themes {
		if theme.Name == name {
			return &theme, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("could not find theme %s", name))
}

func (wp *Wordpress) HasThemeByName(name string) bool {
	_, err := wp.GetThemeByName(name)

	return err == nil
}

func (wp *Wordpress) AddTheme(theme Theme) {
	wp.Themes = append(wp.Themes, theme)
}

func (s *ConfigSite) GetSiteConfigUrl(storeRoot string) string {
	return fmt.Sprintf("%s/sites/%s", storeRoot, path.Join(s.Name, "siteConfig.json"))
}
