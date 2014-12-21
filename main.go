package main

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"time"
	"encoding/json"
	"math/rand"
	"net"
	"fmt"
	"strings"
)

var upgrader = &websocket.Upgrader{}

var H = Hub{
	broadcast: make(chan []byte),
	Register: make(chan *Connection),
	unregister: make(chan *Connection),
	connections: make(map[*Connection]bool),
}

var state = State{
	BallPosition:[2]int{1, 1},
	BallVector: [2]int{1, 1},
	BoardSizeX: 40,
	BoardSizeY: 17,
	PaddleLength: 6,
	PlayerOnePaddle: 0,
	PlayerTwoPaddle: 0,
	PlayerOneVotesUp: make(map[string]bool),
	PlayerOneVotesDown: make(map[string]bool),
	PlayerTwoVotesUp: make(map[string]bool),
	PlayerTwoVotesDown: make(map[string]bool),
	CommandChannel: make(chan IncomingCommand),
}

type Gamedata struct {
	Size [2]int `json:"size"`
	PaddleLength int `json:"paddleLength"`
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
	data := Gamedata{
		Size: [2]int{state.BoardSizeX, state.BoardSizeY},
		PaddleLength: state.PaddleLength,
	}
	jsonData,_ := json.Marshal(&data)
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8000")
	w.Write(jsonData)
}

func main() {
	rand.Seed( time.Now().UTC().UnixNano())

	http.HandleFunc("/state", serveWs)
	http.HandleFunc("/gamedata", serveGamedata)

	go startInputServer()

	go H.run()
	go runGame()
	go state.listenIncomingCommands()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func startInputServer() {
	tcpin,_ := net.Listen("tcp", ":1337")

	for {
		conn,_ := tcpin.Accept()
		go handleInputConnection(conn)
	}
}

type IncomingCommand struct {
	Ip string
	Command string
}

func handleInputConnection(conn net.Conn) {
	defer conn.Close()
	for {
		message := make([]byte, 1024);
		conn.Read(message)

		state.CommandChannel <- IncomingCommand{
			Ip: strings.Split(conn.RemoteAddr().String(), ":")[0],
			Command: fmt.Sprintf("%s", message),
		}
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
