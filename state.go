package main

type State struct {
	BoardSizeX int `json:"-"`
	BoardSizeY int `json:"-"`

	BallPosition [2]int `json:"b"`
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
