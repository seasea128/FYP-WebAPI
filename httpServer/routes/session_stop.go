package routes

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/seasea128/FYP-WebAPI/httpServer/request"
	"github.com/seasea128/FYP-WebAPI/mqttServer"
	"github.com/seasea128/FYP-WebAPI/protobuf/controllerMessage"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

var SessionStopOperation = huma.Operation{
	OperationID: "post-session-stop",
	Method:      http.MethodPost,
	Path:        "/session/stop",
	Summary:     "Set session state to stop",
	Description: "Set session state of a given controller to stop",
	Tags:        []string{"Session"},
}

func SessionStop(db *gorm.DB, mqtt *mqtt.Server) sessionHandler {
	return func(ctx context.Context, input *request.Session) (*request.Session, error) {
		slog.Info("SessionStart", slog.String("body", fmt.Sprintf("%+#v", input.Body)))
		session := &controllerMessage.Session{
			Id:           0,
			ControllerId: input.Body.ControllerID,
			Start:        false,
		}
		sessionOut, err := proto.Marshal(session)
		if err != nil {
			slog.Error("Cannot serialize session MQTT message", slog.String("error", err.Error()))
			return nil, err
		}

		err = mqtt.Publish(mqttServer.SESSION, sessionOut, false, 0)
		if err != nil {
			slog.Error("Cannot send session MQTT message", slog.String("error", err.Error()))
			return nil, err
		}
		return input, nil
	}
}
