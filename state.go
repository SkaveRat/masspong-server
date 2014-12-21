package main

import "math/rand"

type State struct {
	BoardSizeX int `json:"-"`
	BoardSizeY int `json:"-"`

	PaddleLength int `json:"-"`

	BallPosition [2]int `json:"b"`
	BallVector [2]int `json:"-"`

	PlayerOnePaddle int `json:"p1"`
	PlayerTwoPaddle int `json:"p2"`

	PlayerOneScore int `json:"s1"`
	PlayerTwoScore int `json:"s2"`
}

func (s *State) Tick() {
	if (s.isHiddingBoardBorder()) {
		s.reverseY()
	}

	if (s.nextTickIsPaddle()) {
		s.reverseX()
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

	s.PlayerOnePaddle = rand.Intn(s.BoardSizeY-s.PaddleLength)
	s.PlayerTwoPaddle = rand.Intn(s.BoardSizeY-s.PaddleLength)
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

func (s *State) isHiddingBoardBorder() bool {
	return s.BallPosition[1] >= (s.BoardSizeY-1) || s.BallPosition[1] <= 0
}

func (s *State) reverseY() {
	s.BallVector = multiplyVector(s.BallVector, [2]int{1, -1}) //reverse direction
}

func (s *State) moveBall() {
	s.BallPosition = sumVector(s.BallPosition, s.BallVector)
}
