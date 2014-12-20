package main

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"time"
	"encoding/json"
)

var upgrader = &websocket.Upgrader{}

var H = Hub{
	broadcast: make(chan []byte),
	Register: make(chan *Connection),
	unregister: make(chan *Connection),
	connections: make(map[*Connection]bool),
}
type State struct {
	BoardSizeX int
	BoardSizeY int

	BallPosition [2]int `json:"position"`
	BallVector [2]int `json:"-"`
}

func (s *State) Tick() {
	if (s.BallPosition[0] >= (s.BoardSizeX-1) || s.BallPosition[0] <= 0) {
		s.BallVector = multiplyVector(s.BallVector, [2]int{-1, 1}) //reverse direction
	}
	if (s.BallPosition[1] >= (s.BoardSizeY-1) || s.BallPosition[1] <= 0) {
		s.BallVector = multiplyVector(s.BallVector, [2]int{1, -1}) //reverse direction
	}
	s.BallPosition = sumVector(s.BallPosition, s.BallVector)
}


var state = State{
	BallPosition:[2]int{1, 1},
	BallVector: [2]int{1, 1},
	BoardSizeX: 30,
	BoardSizeY: 15,
}


func serveWs(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Origin") //remove same origin policy as dev-hack
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &Connection{send: make(chan []byte), ws: ws}
	H.Register <- c
	go c.writePump()
}

func main() {
	http.HandleFunc("/state", serveWs)
	go H.run()
	go runGame()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}else {
		log.Println("started")
	}
}

func runGame() {

	ticker := time.NewTicker(time.Millisecond * 50)

	for _ = range ticker.C {
		state.Tick()
		value, _ := json.Marshal(&state)
		H.broadcast <- value
	}
}

func multiplyVector(x [2]int, y [2]int) [2]int {
	newVector := [2]int {
			x[0] * y[0],
			x[1] * y[1],
	}
	return newVector
}

func sumVector(x [2]int, y [2]int) [2]int {
	newVector := [2]int{
			x[0] + y[0],
			x[1] + y[1],
	}

	return newVector
}
