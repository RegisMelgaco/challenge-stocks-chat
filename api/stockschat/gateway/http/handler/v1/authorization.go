package v1

import "github.com/regismelgaco/go-sdks/auth/auth/entity"

type Authorization struct {
	Token entity.Token `json:"token"`
}
