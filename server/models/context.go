package models

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type PMContext struct {
	IsReady     bool
	Readability *ReaderConfig
	AWS         *aws.Config
	S3Config    *S3Config
}

func (cfg *PMContext) LoadFromEnv() {
	// Fetch readability config
	readerConfigDir := os.Getenv("PAGEMAIL_READABILITY_CONTEXT_DIR")
	if readerConfigDir == "" {
		panic("readability config directory not set")
	}
	readerPath, err := filepath.Abs(readerConfigDir)
	if err != nil {
		log.Panicf("Could recognise reader config dir given: %s", err)
	}
	_, err = os.Stat(readerPath)
	if err != nil {
		log.Panicf("Could not stat reader context path: %s", err)
	}
	cfg.Readability = &ReaderConfig{
		NodeScript:   "main.js",
		PythonScript: "main.py",
		ContextDir:   readerConfigDir,
	}
	c, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	awsConf, err := config.LoadDefaultConfig(c)
	if err != nil {
		log.Panicf("Could not load AWS config: %s", err)
	}
	cfg.AWS = &awsConf

	cfg.S3Config =&S3Config{
		ReadabilityBucket: "pagemail-speechsynthesis",
	}

	cfg.IsReady = true
}
