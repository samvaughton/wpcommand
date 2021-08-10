package config

import (
	"fmt"
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
		return &types.Wordpress{Plugins: []types.Plugin{}}
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

	if len(config.Wordpress.Sites) > 0 {
		manifest.Sites = config.Wordpress.Sites
	}

	mergePlugins(config, manifest)
	mergeThemes(config, manifest)

	log.Debug("remote manifest file merged with local")

	return manifest
}

func mergePlugins(config *types.Config, manifest *types.Wordpress) {
	for i, mPlugin := range manifest.Plugins {
		if config.Wordpress.HasPluginByName(mPlugin.Name) == false {
			continue
		}

		cPlugin, _ := config.Wordpress.GetPluginByName(mPlugin.Name)

		if cPlugin.Status != "" {
			manifest.Plugins[i].Status = cPlugin.Status
		}

		if cPlugin.Url != "" {
			manifest.Plugins[i].Url = cPlugin.Url
		}

		if cPlugin.Version != "" {
			manifest.Plugins[i].Version = cPlugin.Version
		}

		if cPlugin.Update != "" {
			manifest.Plugins[i].Update = cPlugin.Update
		}
	}

	// now add all plugins that did not exist in the manifest to it
	for _, cPlugin := range config.Wordpress.Plugins {
		// skip if it already has it
		if manifest.HasPluginByName(cPlugin.Name) == true {
			continue
		}

		manifest.AddPlugin(cPlugin)
	}

	// now loop through all plugins and assign their download url
	for i, mPlugin := range manifest.Plugins {
		manifest.Plugins[i].Url = getPluginDownloadUrl(manifest, &mPlugin)
	}
}

func getPluginDownloadUrl(manifest *types.Wordpress, plugin *types.Plugin) string {
	if plugin.Url != "" {
		return plugin.Url
	}

	return fmt.Sprintf("%s/plugins/%s.zip", manifest.StoreRoot, plugin.Name)
}

func mergeThemes(config *types.Config, manifest *types.Wordpress) {
	for i, mTheme := range manifest.Themes {
		if config.Wordpress.HasThemeByName(mTheme.Name) == false {
			continue
		}

		cTheme, _ := config.Wordpress.GetThemeByName(mTheme.Name)

		if cTheme.Status != "" {
			manifest.Themes[i].Status = cTheme.Status
		}

		if cTheme.Url != "" {
			manifest.Themes[i].Url = cTheme.Url
		}

		if cTheme.Version != "" {
			manifest.Themes[i].Version = cTheme.Version
		}

		if cTheme.Update != "" {
			manifest.Themes[i].Update = cTheme.Update
		}
	}

	// now add all themes that did not exist in the manifest to it
	for _, cTheme := range config.Wordpress.Themes {
		// skip if it already has it
		if manifest.HasThemeByName(cTheme.Name) == true {
			continue
		}

		manifest.AddTheme(cTheme)
	}

	// now loop through all themes and assign their download url
	for i, mTheme := range manifest.Themes {
		manifest.Themes[i].Url = getThemeDownloadUrl(manifest, &mTheme)
	}
}

func getThemeDownloadUrl(manifest *types.Wordpress, theme *types.Theme) string {
	if theme.Url != "" {
		return theme.Url
	}

	return fmt.Sprintf("%s/themes/%s.zip", manifest.StoreRoot, theme.Name)
}
