package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/seasea128/FYP-WebAPI/httpServer/request"
)

var GreetingsOperation = huma.Operation{
	OperationID: "get-greeting",
	Method:      http.MethodGet,
	Path:        "/greeting/{name}",
	Summary:     "Get a greeting",
	Description: "Get a greeting for a person by name.",
	Tags:        []string{"Greetings"},
}

func GreetingsHandler(ctx context.Context, input *struct {
	Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
}) (*request.GreetingOutput, error) {
	resp := &request.GreetingOutput{}
	resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
	return resp, nil
}
