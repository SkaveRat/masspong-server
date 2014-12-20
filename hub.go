package main

type Hub struct {
	// Registered connections.
	connections map[*Connection]bool
	// Inbound messages from the connections.
	broadcast chan []byte
	// Register requests from the connections.
	Register chan *Connection
	// Unregister requests from connections.
	unregister chan *Connection
}


func (h *Hub) run() {
	for {
		select {
		case c := <-h.Register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.broadcast:
		for c := range h.connections {
			select {
			case c.send <- m:
			default:
				close(c.send)
				delete(h.connections, c)
			}
		}
		}
	}
}
