package config

import (
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"net/http"
)

func getComputedManifest(config *types.Config) *types.Wordpress {
	return mergeManifestWithLocalConfig(config, downloadAndParseManifest(config))
}

func downloadAndParseManifest(config *types.Config) *types.Wordpress {
	url := config.ManifestUrl

	if url == "" {
		// Return empty
		return &types.Wordpress{}
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("could not retrieve manifest file", err)
	}

	defer resp.Body.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(resp.Body)

	wp := types.Wordpress{}

	// Start YAML decoding from file
	if err := d.Decode(&wp); err != nil {
		log.Fatal(err)
	}

	log.Debug("downloaded and parsed manifest file")

	return &wp
}

func mergeManifestWithLocalConfig(config *types.Config, manifest *types.Wordpress) *types.Wordpress {
	// update the manifest plugins with the local plugins data

	if len(config.Wordpress.DataUrls) > 0 {
		manifest.DataUrls = config.Wordpress.DataUrls
	}

	log.Debug("remote manifest file merged with local")

	return manifest
}
