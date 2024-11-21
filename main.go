package main

import (
	"fmt"

	"github.com/seasea128/WebAPI/config"
	"github.com/seasea128/WebAPI/database"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %s\n", err.Error())
		return
	}

	cfg := config.Load()
	fmt.Printf("Config loaded: %+#v\n", cfg)

	db, err := database.InitConnection(cfg)
	if err != nil {
		fmt.Printf("Error initializing db connection: %s\n", err.Error())
	}
}
