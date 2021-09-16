package config

import (
	"flag"
	"fmt"
	"github.com/evalphobia/logrus_sentry"
	"github.com/kelseyhightower/envconfig"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	"strings"
)

var (
	Config          *types.Config
	ServerAddress   string
	StorageHost     string
	Environment     string
	LoggingLevel    string
	DisableManifest bool
)

func InitConfig(configData string) {
	cfgPath, err := ParseFlags()

	if err != nil {
		log.Fatal(err)
	}

	// Create config structure
	config := &types.Config{}

	d := yaml.NewDecoder(strings.NewReader(configData))

	// if path is provided
	if cfgPath != "" {
		// Open config file
		file, err := os.Open(cfgPath)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		d = yaml.NewDecoder(file)
	}

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		log.Fatal(err)
	}

	// Override config based on environmental variables
	err = envconfig.Process("WPC", config)
	if err != nil {
		log.Fatal(err)
	}

	// Override config based on passed flags
	if StorageHost != "" {
		config.StorageHost = StorageHost
		log.Info("storage-host=", StorageHost)
	}

	if ServerAddress != "" {
		config.ServerAddress = ServerAddress
		log.Info("server-address=", ServerAddress)
	}

	if Environment != "" {
		config.Environment = Environment
		log.Info("environment=", Environment)
	}

	if LoggingLevel != "" {
		config.Logging.Level = LoggingLevel
		log.Info("logging-level=", LoggingLevel)
	}

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	restCfg := &rest.Config{}

	if config.UseInClusterConfig {
		restCfg, err = rest.InClusterConfig()

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("provider.kubernetes: failed to get kubernetes config")
		}

		log.Info("provider.kubernetes: using in-cluster configuration")
	} else {
		restCfg, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("provider.kubernetes: failed to get cmd kubernetes config")
		}
	}

	config.K8RestConfig = restCfg

	log.Info("kube host: ", restCfg.Host)

	if DisableManifest == false {
		config.Wordpress = *getComputedManifest(config)
		log.Info("manifest parsed")
	}

	// add sentry logrus hook

	// sentry logging
	if config.SentryDsn != "" {
		/*err = sentry.Init(sentry.ClientOptions{
			Dsn: config.SentryDsn,
		})

		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}*/

		// the logrus hook initializes it anyway
		hook, err := logrus_sentry.NewSentryHook(config.SentryDsn, []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
		})

		if err == nil {
			log.AddHook(hook)
		}
	}

	Config = config
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	if path == "" {
		return nil
	}

	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}
