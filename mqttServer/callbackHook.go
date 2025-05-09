package mqttServer

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/seasea128/FYP-WebAPI/database/model"
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
	slog.Info("Client published message", slog.String("client", cl.ID),
		slog.String("message", string(pk.Payload[:])),
		slog.String("topic", pk.TopicName))

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
		return
	}

	newSession := &model.Sessions{
		Name: fmt.Sprintf("%s-%d", session.ControllerId, session.SessionId),
	}

	result := h.config.DB.Create(newSession)

	if result.Error != nil {
		slog.Error("Cannot add session to database", slog.String("err", result.Error.Error()))
	}

	slog.Info("Session added to Database", slog.String("controllerId", session.ControllerId), slog.String("session", fmt.Sprintf("%+#v", newSession)))
}

func (h *CallbackHook) handleData(cl *mqtt.Client, pk packets.Packet) {
	data := &controllerMessage.Packet{}

	packetString := string(pk.Payload)

	// TODO: Base64 or hex here?
	//packetBytes, err := hex.DecodeString(packetString)
	//if err != nil {
	//	slog.Error("Cannot convert string to byte array", slog.String("err", err.Error()))
	//	return
	//}

	packetBytes, err := base64.StdEncoding.DecodeString(packetString)
	if err != nil {
		slog.Error("Cannot convert string to byte array", slog.String("err", err.Error()))
		return
	}

	if err := proto.Unmarshal(packetBytes, data); err != nil {
		slog.Error("Cannot deserialize message", slog.String("err", err.Error()))
		return
	}

	switch data.Type {
	case controllerMessage.PacketType_DATA:
		h.handleDataPointsPacket(data.Data)
	case controllerMessage.PacketType_SESSION:
		h.handleSessionPacket(data.Session)
	default:
		h.Log.Error("unexpected controllerMessage.PacketType")
	}

}

func (h *CallbackHook) handleSessionPacket(session *controllerMessage.Session) {
	var controller model.Controllers
	result := h.config.DB.Where(&model.Controllers{ControllerName: session.ControllerId}).First(&controller)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			slog.Error("Cannot find controller with given controller name", slog.String("error", result.Error.Error()), slog.String("controllerName", session.ControllerId))
			return
		}
	}

	if result.RowsAffected == 0 {
		controller.ControllerName = session.ControllerId
		result = h.config.DB.Create(&controller)
		if result.Error != nil {
			slog.Error("Cannot create new controller", slog.String("error", result.Error.Error()))
			return
		}
	}

	if !session.IsActive {
		sessionDB := model.Sessions{
			Name:         fmt.Sprintf("%s-%d", session.ControllerId, session.SessionId),
			ControllerID: session.ControllerId,
			SessionID:    session.SessionId,
		}
		result = h.config.DB.First(&sessionDB)
		if result.Error != nil {
			slog.Error("Cannot query session in database", slog.String("err", result.Error.Error()))
			return
		}

		sessionDB.FinishedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

		return
	}

	sessionDB := model.Sessions{
		ID:           0,
		Name:         fmt.Sprintf("%s-%d", session.ControllerId, session.SessionId),
		ControllerID: session.ControllerId,
		SessionID:    session.SessionId,
		CreatedAt:    time.Time{},
		DeletedAt:    sql.NullTime{},
	}

	result = h.config.DB.Create(&sessionDB)

	if result.Error != nil {
		slog.Error("Cannot add session to database", slog.String("err", result.Error.Error()))
		return
	}

	slog.Info("Session added to Database", slog.String("controllerId", session.ControllerId), slog.String("session", fmt.Sprintf("%+#v", sessionDB)))
}

func (h *CallbackHook) handleDataPointsPacket(data *controllerMessage.DataPoints) {
	for i, measurement := range data.Measurement {
		suspensionLog := model.SuspensionLogs{
			CreatedAt:    data.Timestamp.AsTime().Add(time.Millisecond * time.Duration(10*i)),
			ControllerID: data.ControllerId,
			SessionID:    data.SessionId,
			LeftTop:      measurement.DistanceLt,
			LeftBottom:   measurement.DistanceLb,
			RightTop:     measurement.DistanceRt,
			RightBottom:  measurement.DistanceRb,
			GPSPosition:  "",
			GPSSpeed:     "",
		}

		result := h.config.DB.Create(&suspensionLog)

		if result.Error != nil {
			slog.Error("Cannot add data to database", slog.String("err", result.Error.Error()))
		}

		slog.Info("Data added to Database", slog.String("controllerId", data.ControllerId), slog.String("data", fmt.Sprintf("%+#v", suspensionLog)))

		h.config.Server.Publish(fmt.Sprintf("realtime/%s/DistanceLt", data.ControllerId), []byte(strconv.Itoa(int(measurement.DistanceLt))), false, 0)
		h.config.Server.Publish(fmt.Sprintf("realtime/%s/DistanceLb", data.ControllerId), []byte(strconv.Itoa(int(measurement.DistanceLb))), false, 0)
		h.config.Server.Publish(fmt.Sprintf("realtime/%s/DistanceRt", data.ControllerId), []byte(strconv.Itoa(int(measurement.DistanceRt))), false, 0)
		h.config.Server.Publish(fmt.Sprintf("realtime/%s/DistanceRb", data.ControllerId), []byte(strconv.Itoa(int(measurement.DistanceRb))), false, 0)

	}
}
