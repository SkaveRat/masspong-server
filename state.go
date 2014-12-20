package main

type State struct {
	BoardSizeX int `json:"-"`
	BoardSizeY int `json:"-"`

	PaddleLength int `json:"-"`

	BallPosition [2]int `json:"b"`
	BallVector [2]int `json:"-"`

	PlayerOnePaddle int `json:"p1"`
	PlayerTwoPaddle int `json:"p2"`
}

func (s *State) Tick() {
	if (s.isHiddingBoardBorder()) {
		s.reverseY()
	}

	if(s.nextTickIsPaddle()) {
		s.reverseX()
	}

	if (s.isOverPlayerBorder()) {
		s.Reset()
	}
	s.moveBall()
}

func (s *State) nextTickIsPaddle() bool {
	nextStep := sumVector(s.BallPosition, s.BallVector)

	paddleLength := 5

	if(s.BallVector[0] > 0 && (nextStep[0] >= s.BoardSizeX-1)) { //moving right && is paddle-column
		return (nextStep[1] >= s.PlayerTwoPaddle && nextStep[1] <= s.PlayerTwoPaddle + paddleLength)
	}else if(s.BallVector[0] < 0 && (nextStep[0] <= 0)){ // moving left && is paddle column
		return (nextStep[1] >= s.PlayerOnePaddle && nextStep[1] <= s.PlayerOnePaddle + paddleLength)
	}else{
		return false
	}
}

func (s *State) Reset() {
	s.BallPosition = [2]int{3,3}
	s.BallVector = [2]int{1,1}
}

func (s *State) isOverPlayerBorder() bool{
	return s.BallPosition[0] >= (s.BoardSizeX-1) || s.BallPosition[0] <= 0
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
