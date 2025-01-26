package domain

import "github.com/gorilla/websocket"

type Client struct {
	User User
	Conn *websocket.Conn
}
