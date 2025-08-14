package config

import (
	"fmt"
	"os"
)

// FindConfigFile ищет конфиги
func FindConfigFile(paths ...string) string {
	if len(paths) == 0 {
		paths = []string{
			"./configs/config.yaml",
			"./go-news/configs/config.yaml",
			"/app/configs/config.yaml",
		}
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("Found config at: %s\n", path)
			return path
		}
		fmt.Printf("Config not found at: %s\n", path)
	}

	fmt.Println("Config file not found in any expected location")

	return ""
}
