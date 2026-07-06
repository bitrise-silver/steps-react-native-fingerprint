package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-steputils/tools"
)

// Config maps the step inputs (see step.yml) to Go fields.
type Config struct {
	FilePaths string `env:"file_paths,required"`
	KeyPrefix string `env:"key_prefix"`
	Verbose   bool   `env:"verbose"`
}

func main() {
	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Invalid input: %s", err)
	}
	stepconf.Print(cfg)
	fmt.Println()

	paths := parsePaths(cfg.FilePaths)
	if cfg.Verbose {
		log("Fingerprinting %d file(s):", len(paths))
		for _, p := range paths {
			log("  - %s", p)
		}
	}

	fingerprint, err := computeFingerprint(paths)
	if err != nil {
		failf("Failed to compute fingerprint: %s", err)
	}
	hashString := withKeyPrefix(cfg.KeyPrefix, fingerprint)

	if err := tools.ExportEnvironmentWithEnvman("BUNDLE_HASH_STRING", hashString); err != nil {
		failf("Failed to export BUNDLE_HASH_STRING: %s", err)
	}

	apiBaseURL := os.Getenv("BITRISEIO_ABCS_API_URL")
	accessToken := os.Getenv("BITRISEIO_BITRISE_SERVICES_ACCESS_TOKEN")
	if apiBaseURL == "" || accessToken == "" {
		failf("Bitrise key-value cache service is not available: BITRISEIO_ABCS_API_URL / BITRISEIO_BITRISE_SERVICES_ACCESS_TOKEN is not set. This step must run on Bitrise CI.")
	}

	found, err := cacheKeyExists(apiBaseURL, accessToken, hashString)
	if err != nil {
		failf("Failed to check cache: %s", err)
	}
	if err := tools.ExportEnvironmentWithEnvman("BUNDLE_CACHE_FOUND", strconv.FormatBool(found)); err != nil {
		failf("Failed to export BUNDLE_CACHE_FOUND: %s", err)
	}

	fmt.Println()
	log("Fingerprint : BUNDLE_HASH_STRING=%s", hashString)
	log("Cache found : BUNDLE_CACHE_FOUND=%t", found)
	os.Exit(0)
}
