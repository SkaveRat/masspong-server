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
	BoardSizeX int `json:"-"`
	BoardSizeY int `json:"-"`

	BallPosition [2]int `json:"b"`
	BallVector [2]int `json:"-"`
}

type Gamedata struct {
	Size [2]int `json:"size"`
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
	BoardSizeX: 40,
	BoardSizeY: 17,
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

func serveGamedata(w http.ResponseWriter, r *http.Request) {
	data := Gamedata{Size: [2]int{state.BoardSizeX, state.BoardSizeY}}
	jsonData,_ := json.Marshal(&data)
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8000")
	w.Write(jsonData)
}

func main() {
	http.HandleFunc("/state", serveWs)
	http.HandleFunc("/gamedata", serveGamedata)
	go H.run()
	go runGame()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func runGame() {

	ticker := time.NewTicker(time.Millisecond * 100)

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
