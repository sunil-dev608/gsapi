package config

import (
	"os"
	"sort"
	"strconv"
	"strings"
)

type Config struct {
	PackSizes []int
}

var cfg *Config

var defaultPackSizes = []int{5000, 2000, 1000, 500, 250}

func GetConfig() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{}

	packSizesEnv := os.Getenv("PACK_SIZES")
	if packSizesEnv == "" {
		cfg.PackSizes = defaultPackSizes
		return cfg, nil
	}
	packSizes := strings.Split(packSizesEnv, ",")
	for _, packSize := range packSizes {
		if p, err := strconv.Atoi(packSize); err != nil {
			return nil, err
		} else {
			cfg.PackSizes = append(cfg.PackSizes, p)
		}
	}

	sort.Slice(cfg.PackSizes, func(i, j int) bool {
		return cfg.PackSizes[i] > cfg.PackSizes[j]
	})

	return cfg, nil
}
