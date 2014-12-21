package main

import (
	"math/rand"
	"regexp"
	"strings"
)

type State struct {
	BoardSizeX int `json:"-"`
	BoardSizeY int `json:"-"`

	PaddleLength int `json:"-"`

	BallPosition [2]int `json:"b"`
	BallVector   [2]int `json:"-"`

	PlayerOnePaddle int `json:"p1"`
	PlayerTwoPaddle int `json:"p2"`

	PlayerOneScore int `json:"s1"`
	PlayerTwoScore int `json:"s2"`

	PlayerOneVotesUp   map[string]bool `json:"-"`
	PlayerOneVotesDown map[string]bool `json:"-"`
	PlayerTwoVotesUp   map[string]bool `json:"-"`
	PlayerTwoVotesDown map[string]bool `json:"-"`

	CommandChannel chan IncomingCommand `json:"-"`
}

func (s *State) Tick() {

	s.movePaddles()

	if (s.nextTickIsPaddle()) {
		s.reverseX()
	}

	if (s.isHittingBoardBorder()) {
		s.reverseY()
	}

	if (s.isOverPlayerOneBorder()) {
		s.PlayerTwoScore++
		s.Reset(1)
	}

	if (s.isOverPlayerTwoBorder()) {
		s.PlayerOneScore++
		s.Reset(-1)
	}

	s.moveBall()
}

func (s *State) movePaddles() {
	playerOneDiff := len(s.PlayerOneVotesDown) - len(s.PlayerOneVotesUp)
	playerTwoDiff := len(s.PlayerTwoVotesDown) - len(s.PlayerTwoVotesUp)


	switch {
	case playerOneDiff < 0:
		s.PlayerOnePaddle--
	case playerOneDiff == 0:
	case playerOneDiff > 0:
		s.PlayerOnePaddle++
	}

	switch {
	case playerTwoDiff < 0:
		s.PlayerTwoPaddle--
	case playerTwoDiff == 0:
	case playerTwoDiff > 0:
		s.PlayerTwoPaddle++
	}

	if(s.PlayerOnePaddle < 0) {s.PlayerOnePaddle = 0}
	if(s.PlayerOnePaddle > s.BoardSizeY - s.PaddleLength) {s.PlayerOnePaddle = s.BoardSizeY - s.PaddleLength}

	if(s.PlayerTwoPaddle < 0) {s.PlayerTwoPaddle = 0}
	if(s.PlayerTwoPaddle > s.BoardSizeY - s.PaddleLength) {s.PlayerTwoPaddle = s.BoardSizeY - s.PaddleLength}


	s.PlayerOneVotesUp   = make(map[string]bool)
	s.PlayerOneVotesDown = make(map[string]bool)
	s.PlayerTwoVotesUp   = make(map[string]bool)
	s.PlayerTwoVotesDown = make(map[string]bool)
}

func (s *State) nextTickIsPaddle() bool {
	nextStep := sumVector(s.BallPosition, s.BallVector)

	if (s.BallVector[0] > 0 && (nextStep[0] >= s.BoardSizeX-1)) { //moving right && is paddle-column
		paddleDiff := nextStep[1] - s.PlayerTwoPaddle
		switch paddleDiff {
		case 0: s.BallVector[1] = -1; break;
		case 1: s.BallVector[1] = -1; break;
		case 2: s.BallVector[1] = 0; break;
		case 3: s.BallVector[1] = 0; break;
		case 4: s.BallVector[1] = 1; break;
		case 5: s.BallVector[1] = 1; break;
		}
		return (paddleDiff >= 0 && paddleDiff < s.PaddleLength)
	}else if (s.BallVector[0] < 0 && (nextStep[0] <= 0)) { // moving left && is paddle column
		paddleDiff := nextStep[1] - s.PlayerOnePaddle
		switch paddleDiff {
		case 0: s.BallVector[1] = -1; break;
		case 1: s.BallVector[1] = -1; break;
		case 2: s.BallVector[1] = 0; break;
		case 3: s.BallVector[1] = 0; break;
		case 4: s.BallVector[1] = 1; break;
		case 5: s.BallVector[1] = 1; break;
		}
		return (paddleDiff >= 0 && paddleDiff < s.PaddleLength)
	}else {
		return false
	}
}

func (s *State) Reset(initialDirection int) {
	var initialX int;
	if (initialDirection < 0 ) {
		initialX = s.BoardSizeX-3
	}else {
		initialX = 3
	}
	s.BallPosition = [2]int{initialX, rand.Intn(s.BoardSizeY-1)}
	s.BallVector = [2]int{initialDirection, 1}
}

func (s *State) isOverPlayerOneBorder() bool {
	return s.BallPosition[0] <= 0
}

func (s *State) isOverPlayerTwoBorder() bool {
	return s.BallPosition[0] >= (s.BoardSizeX-1)
}

func (s *State) reverseX() {
	s.BallVector = multiplyVector(s.BallVector, [2]int{-1, 1}) //reverse direction
}

func (s *State) isHittingBoardBorder() bool {
	return s.BallPosition[1] >= (s.BoardSizeY-1) || s.BallPosition[1] <= 0
}

func (s *State) reverseY() {
	s.BallVector = multiplyVector(s.BallVector, [2]int{1, -1}) //reverse direction
}

func (s *State) moveBall() {
	s.BallPosition = sumVector(s.BallPosition, s.BallVector)
}

func (s *State) listenIncomingCommands() {
	for command := range s.CommandChannel {
		matched,_ := regexp.Match(`^[1|2] (UP|DOWN|up|down)`, []byte(command.Command))
		if(matched) {

			cmd := strings.Replace(command.Command, "\n", "", 1)

			splitted := strings.Split(cmd, " ")

			matchDirection,_ := regexp.Match(`(up)`, []byte(splitted[1]))
			if(splitted[0] == "1") {
				if (matchDirection) {
					s.PlayerOneVotesUp[command.Ip] = true
				}else{
					s.PlayerOneVotesDown[command.Ip] = true
				}
			}else{
				if(matchDirection) {
					s.PlayerTwoVotesUp[command.Ip] = true
				}else{
					s.PlayerTwoVotesDown[command.Ip] = true
				}
			}
		}
	}
}
