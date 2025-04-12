package config

import (
	"flag"
	"os"
)

// Constants for audio configuration
const (
	SampleRate      = 16000
	FramesPerBuffer = 1024
	Channels        = 1
)

// AppConfig holds the application-wide configuration
type AppConfig struct {
	ModelPath      string
	LogFilePath    string
}

// NewConfig creates and initializes a new configuration
func NewConfig() *AppConfig {
	cfg := &AppConfig{
		LogFilePath:    "speech-reco.log",
	}

	// Parse command line flags
	flag.StringVar(&cfg.ModelPath, "model", "models/ggml-base.en.bin", "Path to Whisper model file")
	flag.Parse()

	// Validate model path
	if _, err := os.Stat(cfg.ModelPath); os.IsNotExist(err) {
		// Will be handled by the caller
		return cfg
	}

	return cfg
}