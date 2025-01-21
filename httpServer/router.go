package httpServer

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"gorm.io/gorm"
)

type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func CreateHTTPServer(db *gorm.DB) *http.ServeMux {
	router := http.NewServeMux()

	api := humago.New(router, huma.DefaultConfig("Testing", "1.0.0"))

	api.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		slog.Info(ctx.URL().Path)
	})

	AddRoutes(api)

	return router
}

func AddRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-greeting",
		Method:      http.MethodGet,
		Path:        "/greeting/{name}",
		Summary:     "Get a greeting",
		Description: "Get a greeting for a person by name.",
		Tags:        []string{"Greetings"},
	}, func(ctx context.Context, input *struct {
		Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
	}) (*GreetingOutput, error) {
		resp := &GreetingOutput{}
		resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
		return resp, nil
	})
}
