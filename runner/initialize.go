package runner

import (
	"strings"

	"iplookup/passive"
)

// initializePassiveEngine creates the passive engine and loads sources etc
func (r *Runner) initializePassiveEngine() {
	var sources, exclusions []string

	if r.options.ExcludeSources != "" {
		exclusions = append(exclusions, strings.Split(r.options.ExcludeSources, ",")...)
	} else {
		exclusions = append(exclusions, r.options.YAMLConfig.ExcludeSources...)
	}

	// Use all sources if asked by the user
	if r.options.All {
		sources = append(sources, r.options.YAMLConfig.AllSources...)
	}

	// If there are any sources from CLI, only use them
	// Otherwise, use the yaml file sources
	if !r.options.All {
		if r.options.Sources != "" {
			sources = append(sources, strings.Split(r.options.Sources, ",")...)
		} else {
			sources = append(sources, r.options.YAMLConfig.Sources...)
		}
	}
	r.passiveAgent = passive.New(sources, exclusions)
}
