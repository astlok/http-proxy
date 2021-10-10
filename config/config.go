package config

import (
	"fmt"
	"os"
)

type config struct {
	HTTPS bool
}

func NewConfig() *config {
	var https bool
	file, err := os.Open("config/conf")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	data := make([]byte, 64)

	n, err := file.Read(data)

	switch string(data[:n]) {
	case "http":
		https = false
	case "https":
		https = true

	}
	return &config{
		HTTPS: https,
	}
}
