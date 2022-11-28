package main

import (
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type SourceDependenciesMap map[string][]*Dependency

type Dependency struct {
	Alias *string
	SourceRef
}

func (d *Dependency) GetTargetPath(basePath string) (string, error) {
	var relativePath string
	if d.Alias == nil {
		// Use segment target path format
		pathParts := strings.Split(d.Source, ":")
		repositoryName := pathParts[1]
		relativePath = filepath.Join(basePath, repositoryName, d.Version)
	} else {
		// Use community target path format
		relativePath = filepath.Join(basePath, *d.Alias)
	}

	return filepath.Abs(relativePath)
}

type SourceRef struct {
	Source  string `yaml:"source"`
	Version string `yaml:"version"`
}

func parseTerrafile(in []byte) (SourceDependenciesMap, error) {
	// Try parse Segment internal format
	result, err := parseSegmentTerrafile(in)
	if _, ok := err.(*yaml.TypeError); ok {
		// Try fallback to community format
		result, err = parseCommunityTerrafile(in)
	}

	return result, err
}

func parseSegmentTerrafile(in []byte) (SourceDependenciesMap, error) {
	var config map[string][]string
	if err := yaml.Unmarshal(in, &config); err != nil {
		return nil, err
	}

	result := make(SourceDependenciesMap)
	for source, versions := range config {
		for _, version := range versions {
			result[source] = append(result[source], &Dependency{
				SourceRef: SourceRef{
					Source:  source,
					Version: version,
				},
			})
		}
	}

	return result, nil
}

func parseCommunityTerrafile(in []byte) (SourceDependenciesMap, error) {
	var config map[string]SourceRef
	if err := yaml.Unmarshal(in, &config); err != nil {
		return nil, err
	}

	result := make(SourceDependenciesMap)
	for key, sourceRef := range config {
		alias := key
		result[sourceRef.Source] = append(result[sourceRef.Source], &Dependency{
			Alias:     &alias,
			SourceRef: sourceRef,
		})
	}

	return result, nil
}
