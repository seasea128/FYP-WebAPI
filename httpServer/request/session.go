package request

type Session struct {
	Body sessionBody `contentType:"application/json" required:"true"`
}

type sessionBody struct {
	ControllerID string `doc:"ID of controller to set state"`
}
