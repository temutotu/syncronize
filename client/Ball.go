package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type Ball struct {
	x         int
	y         int
	xVelocity int
	yVelocity int
	timer     time.Time
	// ボールが反射する境界情報
	top       int
	bottom    int
	leftWall  int
	rightWall int
}

func NewBall(x int, y int, height int, width int) *Ball {
	start := time.Now()
	ball := Ball{x, y, 1, 1, start, 3, height - 2, 3, width - 4}
	return &ball
}

func (ball *Ball) MoveBall() {
	xMoved := ball.x + ball.xVelocity
	yMoved := ball.y + ball.yVelocity
	if ball.CheckCollision(xMoved, yMoved) {
		xMoved = ball.x + ball.xVelocity
		yMoved = ball.y + ball.yVelocity
	}
	ball.x = xMoved
	ball.y = yMoved
}

func (ball *Ball) Update(screen *tcell.Screen) {
	elapsed := time.Since(ball.timer)
	if elapsed < 100*time.Millisecond {
		return
	}
	ball.timer = time.Now()
	DrawDeleteBall(*screen, *ball)
	ball.MoveBall()
	//DrawBall(*screen, *ball)
}

func (ball *Ball) CheckCollision(xMoved int, yMoved int) bool {
	// y方向衝突判定
	if yMoved < ball.top || yMoved > ball.bottom {
		ball.yVelocity = ball.yVelocity * (-1)
		return true
	}
	// x方向衝突判定
	if xMoved < ball.leftWall || xMoved > ball.rightWall {
		ball.xVelocity = ball.xVelocity * (-1)
		return true
	}
	return false
}
