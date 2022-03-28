package api

type EventBody struct {
	Message string `json:"message"`
}

type Error struct {
	Detail string `json:"detail"`
}
