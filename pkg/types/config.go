package types

import (
	"k8s.io/client-go/rest"
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
	StoreRoot string    `yaml:"storeRoot" envconfig:"WP_STORE_ROOT" json:"storeRoot"`
	DataUrls  []DataUrl `yaml:"dataUrls" json:"dataUrls"`
}

type Config struct {
	K8RestConfig *rest.Config `json:"-"`
	// Runtime not meant for env or yaml files, only flags
	Runtime struct {
		SpecificPod string
	} `json:"-"`

	StorageHost   string `yaml:"storageHost" envconfig:"STORAGE_HOST" json:"storageHost"`
	ServerAddress string `yaml:"serverAddress" envconfig:"SERVER_ADDRESS" json:"serverAddress"`
	Environment   string `yaml:"environment" envconfig:"ENVIRONMENT" json:"environment"`
	DatabaseDsn   string `yaml:"databaseDsn" envconfig:"DATABASE_DSN" json:"-"`

	K8 struct {
		LabelSelector string `yaml:"labelSelector" envconfig:"K8_LABEL_SELECTOR" json:"labelSelector"`
	} `yaml:"kubernetes" json:"kubernetes"`

	Logging struct {
		Level string `yaml:"level" envconfig:"LOGGING_LEVEL" json:"level"`
	} `yaml:"logging" json:"logging"`

	ManifestUrl string    `yaml:"manifest" envconfig:"WP_MANIFEST" json:"manifest"`
	Wordpress   Wordpress `yaml:"wordpress" json:"wordpress"`
}
