package mqttServer

import (
	"bytes"
	"log/slog"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/seasea128/FYP-WebAPI/protobuf/controllerMessage"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

type CallbackHook struct {
	mqtt.HookBase
	config *CallbackHookOptions
}

type CallbackHookOptions struct {
	Server *mqtt.Server
	DB     *gorm.DB
}

func (h *CallbackHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnDisconnect,
		mqtt.OnSubscribed,
		mqtt.OnUnsubscribed,
		mqtt.OnPublished,
	}, []byte{b})
}

func (h *CallbackHook) Init(config any) error {
	if _, ok := config.(*CallbackHookOptions); !ok && config != nil {
		return mqtt.ErrInvalidConfigType
	}

	h.config = config.(*CallbackHookOptions)
	if h.config.Server == nil {
		return mqtt.ErrInvalidConfigType
	}
	if h.config.DB == nil {
		return mqtt.ErrInvalidConfigType
	}
	return nil
}

func (h *CallbackHook) OnConnect(cl *mqtt.Client, pk packets.Packet) error {
	slog.Info("Client connected", slog.String("client", cl.ID))
	return nil
}

func (h *CallbackHook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	var errStr string
	if err != nil {
		errStr = err.Error()
	} else {
		errStr = "nil"
	}
	slog.Info("Client disconnected", slog.String("client", cl.ID),
		slog.String("error", errStr),
		slog.Bool("expires", expire))
}

func (h *CallbackHook) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	slog.Info("Client subscribed", slog.String("client", cl.ID), slog.String("topic", pk.TopicName))
}

func (h *CallbackHook) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {
	slog.Info("Client unsubscribed", slog.String("client", cl.ID), slog.String("topic", pk.TopicName))
}

func (h *CallbackHook) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	slog.Info("Client published message", slog.String("client", cl.ID), slog.String("message", string(pk.Payload[:])), slog.String("topic", pk.TopicName))

	switch pk.TopicName {
	case DATA:
		{
			h.handleData(cl, pk)
		}
	case SESSION:
		{
			h.handleSession(cl, pk)
		}
	}
}

func (h *CallbackHook) handleSession(cl *mqtt.Client, pk packets.Packet) {
	session := &controllerMessage.Session{}

	if err := proto.Unmarshal(pk.Payload, session); err != nil {
		slog.Error("Cannot deserialize message", slog.String("err", err.Error()))
	}
}

func (h *CallbackHook) handleData(cl *mqtt.Client, pk packets.Packet) {
	data := &controllerMessage.DataReceived{}

	if err := proto.Unmarshal(pk.Payload, data); err != nil {
		slog.Error("Cannot deserialize message", slog.String("err", err.Error()))
	}

}
