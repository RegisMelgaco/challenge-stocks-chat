package v1

type ErrorOutput struct {
	Message string `json:"message"`
}

func ToErrorOutput(msg string) Response {
	return Response{
		Type:    "error",
		Payload: ErrorOutput{Message: msg},
	}
}
