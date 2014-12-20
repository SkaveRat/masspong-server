package main

import (
	"golang.org/x/net/websocket"
	"net/http"
	"time"
	"encoding/json"
)

type State struct {
	BallPosition [2] int `json:"position"`
	BallVector [2]int `json:"-"`
}

func (s *State) Tick() {
	if(s.BallPosition[0] >= 19 || s.BallPosition[0] <= 0) {
		s.BallVector = multiplyVector(s.BallVector, [2]int{-1,1}) //reverse direction
	}
	if(s.BallPosition[1] >= 9 || s.BallPosition[1] <= 0) {
		s.BallVector = multiplyVector(s.BallVector, [2]int{1,-1}) //reverse direction
	}
	s.BallPosition = sumVector(s.BallPosition, s.BallVector)
}

func EchoServer(ws *websocket.Conn) {

	state := State{}
	state.BallPosition = [2]int{1,1}
	state.BallVector = [2]int{1,1}
	timer := time.NewTicker(time.Millisecond * 50)

	for _ = range timer.C {
		state.Tick()
		value,_ := json.Marshal(&state)
		ws.Write(value)
	}
}

func main() {
	http.Handle("/echo", websocket.Handler(EchoServer))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func multiplyVector(x [2]int, y [2]int) [2]int {
	newVector := [2]int{
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
