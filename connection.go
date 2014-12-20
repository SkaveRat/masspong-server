package main

import "github.com/gorilla/websocket"

// connection is an middleman between the websocket connection and the hub.
type Connection struct {
	// The websocket connection.
	ws *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

// write writes a message with the given message type and payload.
func (c *Connection) write(mt int, payload []byte) error {
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Connection) writePump() {
	defer func() {
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}
