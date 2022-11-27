package chat

import (
	"encoding/json"
	v1 "local/challengestockschat/stockschat/gateway/http/handler/chat/v1"
	"net/http"

	"github.com/regismelgaco/go-sdks/erring"
	"github.com/regismelgaco/go-sdks/httpresp"
)

func (h Handler) PostMessage(r *http.Request) httpresp.Response {
	var input v1.InputMessage
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		err = erring.Wrap(err).Build()

		return httpresp.BadRequest(err)
	}

	msg, err := input.ToEntity(r.Context())
	if err != nil {
		return httpresp.Internal(err)
	}

	if err := h.u.CreateMessage(r.Context(), &msg); err != nil {
		return httpresp.Error(err)
	}

	return httpresp.Created(v1.ToMessangeOutput(msg))
}
