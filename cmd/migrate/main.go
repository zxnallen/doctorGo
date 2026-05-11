package main

import (
	"flag"
	"log"

	"doctor-go/internal/config"
	"doctor-go/internal/migration"
)

func main() {
	action := flag.String("action", "up", "migration action: up, down, version")
	steps := flag.Int("steps", 1, "steps for down action")
	version := flag.Int("version", 0, "version for force action")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	switch *action {
	case "up":
		err = migration.Up(cfg.MySQL)
	case "down":
		err = migration.Down(cfg.MySQL, *steps)
	case "version":
		err = migration.PrintVersion(cfg.MySQL)
	case "force":
		err = migration.Force(cfg.MySQL, *version)
	default:
		log.Fatalf("unknown action: %s", *action)
	}
	if err != nil {
		log.Fatalf("migration %s: %v", *action, err)
	}
}
