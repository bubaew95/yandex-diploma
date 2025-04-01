package resplogindto

import "time"

type ResponseToken struct {
	Token  string    `json:"token"`
	Expire time.Time `json:"expire"`
}
