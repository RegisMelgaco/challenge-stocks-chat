package chat

import "time"

type Message struct {
	Author    string
	Content   string
	CreatedAt time.Time
}
