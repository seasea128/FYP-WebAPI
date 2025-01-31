package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/seasea128/FYP-WebAPI/config"
	"github.com/seasea128/FYP-WebAPI/database"
	"github.com/seasea128/FYP-WebAPI/httpServer"
	"github.com/seasea128/FYP-WebAPI/mqttServer"
)

type Options struct {
}

var (
	hash    string
	version string
)

func main() {
	slog.Info("WheelSensor API/MQTT Server")
	slog.Info(fmt.Sprintf("Version: %s", version))
	slog.Info(fmt.Sprintf("Hash: %s", hash))
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %s\n", err.Error())
		return
	}

	cfg := config.Load()
	slog.Info("Config loaded:", slog.String("cfg", fmt.Sprintf("%+#v", cfg)))

	db, err := database.InitConnection(cfg)
	if err != nil {
		fmt.Printf("Error initializing db connection: %s\n", err.Error())
	}

	cli := humacli.New(func(hooks humacli.Hooks, option *Options) {
		mqttHandler, err := mqttServer.CreateMQTTServer(cfg, db)
		if err != nil {
			slog.Error("Failed to create MQTT server", slog.String("error", err.Error()))
			return
		}

		httpHandler := httpServer.CreateHTTPServer(db, mqttHandler)

		hooks.OnStart(func() {
			slog.Info("Starting MQTT server", slog.Int("Port", cfg.MQTTPort))
			err := mqttHandler.Serve()
			if err != nil {
				slog.Error("Failed to start MQTT server", slog.String("error", err.Error()))
				return
			}

			slog.Info("Starting HTTP server", slog.Int("Port", cfg.HTTPPort))
			http.ListenAndServe(fmt.Sprintf(":%d", cfg.HTTPPort), httpHandler)
		})

		hooks.OnStop(func() {
			// TODO: Cleanup
		})
	})

	cli.Run()
}
