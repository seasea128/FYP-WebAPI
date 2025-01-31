package httpServer

import (
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/seasea128/FYP-WebAPI/httpServer/routes"
	"gorm.io/gorm"
)

func CreateHTTPServer(db *gorm.DB, mqtt *mqtt.Server) *http.ServeMux {
	router := http.NewServeMux()

	api := humago.New(router, huma.DefaultConfig("Testing", "1.0.0"))

	// TODO: Request/Response Logger
	api.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		slog.Info(ctx.URL().Path)
		next(ctx)
	})

	// TODO: Token checker middleware

	AddRoutes(api, mqtt, db)

	return router
}

func AddRoutes(api huma.API, mqtt *mqtt.Server, db *gorm.DB) {
	huma.Register(api, routes.GreetingsOperation, routes.GreetingsHandler)
	huma.Register(api, routes.SessionStartOperation, routes.SessionStart(db, mqtt))
	huma.Register(api, routes.SessionStopOperation, routes.SessionStop(db, mqtt))
}
