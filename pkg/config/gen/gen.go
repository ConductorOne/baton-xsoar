package main

import (
	cfg "github.com/conductorone/baton-xsoar/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("xsoar", cfg.Config)
}
