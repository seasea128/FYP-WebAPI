package mqttServer

import (
	"fmt"
	"log/slog"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/seasea128/FYP-WebAPI/config"
	"gorm.io/gorm"
)

func CreateMQTTServer(cfg *config.Configuration, db *gorm.DB) (*mqtt.Server, error) {
	server := mqtt.New(&mqtt.Options{InlineClient: true})
	server.Log = slog.Default()

	err := server.AddHook(new(auth.AllowHook), nil)

	callbackHook := new(CallbackHook)

	// Add hooks for message
	err = server.AddHook(callbackHook, &CallbackHookOptions{Server: server, DB: db})
	if err != nil {
		err = fmt.Errorf("Cannot add hooks: %s", err.Error())
		return nil, err
	}
	tcp := listeners.NewTCP(listeners.Config{ID: "t1", Address: fmt.Sprintf(":%d", cfg.MQTTPort)})
	stats := listeners.NewHTTPStats(listeners.Config{ID: "s1", Address: ":8902"}, server.Info)
	err = server.AddListener(tcp)
	err = server.AddListener(stats)
	if err != nil {
		err = fmt.Errorf("Cannot add listener: %s", err.Error())
		return nil, err
	}

	return server, nil
}
