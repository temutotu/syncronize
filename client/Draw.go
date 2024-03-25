package main

import (
	"github.com/gdamore/tcell/v2"
)

var drawManager *DrawManager = nil

const GAME_SCREEN_TOP int = 3
const GAME_SCREEN_BOTTOM int = 43
const GAME_SCREEN_LEFT int = 1
const GAME_SCREEN_RIGHT int = 158

type DrawManager struct {
	screen *tcell.Screen
}

func InitDrawManager(screen *tcell.Screen) {
	drawManager = &DrawManager{screen}
}

func (manager *DrawManager) NewDrawLog(logYPos int, logMessage string) {
	for i := 0; i < 100; i++ {
		(*manager.screen).SetContent(i, logYPos, ' ', nil, tcell.StyleDefault)
	}

	text := []rune(logMessage)
	for i := 0; i < len(logMessage); i++ {
		(*manager.screen).SetContent(i, logYPos, text[i], nil, tcell.StyleDefault)
	}
	(*manager.screen).Show()
}

func DrawLog(screen tcell.Screen, logYPos int, logMessage string) {
	text := []rune(logMessage)
	for i := 0; i < len(logMessage); i++ {
		screen.SetContent(i, logYPos, text[i], nil, tcell.StyleDefault)
	}
	screen.Show()
}

func DrawError(screen tcell.Screen, errorMessage string) {
	screen.Clear()
	text := []rune(errorMessage)
	for i := 0; i < len(errorMessage); i++ {
		screen.SetContent(i, i, text[i], nil, tcell.StyleDefault)
	}
	screen.Show()
}

func DrawFrame(screen tcell.Screen, height int, width int) {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {

			if j == 0 || j == width-1 {
				screen.SetContent(j, i, '|', nil, tcell.StyleDefault)
			}

			if i == 0 || i == height-1 {
				screen.SetContent(j, i, '=', nil, tcell.StyleDefault)
			}

			if i == 2 {
				if j > 0 && j < width-1 {
					screen.SetContent(j, i, '=', nil, tcell.StyleDefault)
				}
			}
		}

		if i == 1 {
			text := []rune("Esc:quit ↑:UpSlideBar ↓:DownSlideBar")
			for l := 0; l < len(text); l++ {
				screen.SetContent(1+l, i, text[l], nil, tcell.StyleDefault)
			}
		}
	}
	screen.Show()
}

func (manager *DrawManager) CreanGameScreen() {
	for i := GAME_SCREEN_TOP; i <= GAME_SCREEN_BOTTOM; i++ {
		for j := GAME_SCREEN_LEFT; j <= GAME_SCREEN_RIGHT; j++ {
			(*manager.screen).SetContent(j, i, ' ', nil, tcell.StyleDefault)
		}
	}
	(*manager.screen).Show()
}

func (manager *DrawManager) DrawSideBar(screen tcell.Screen, sideBar SideBar) {
	for i := 0; i < sideBar.size; i++ {
		(*manager.screen).SetContent(sideBar.x, sideBar.y+i, '|', nil, tcell.StyleDefault)
		(*manager.screen).SetContent(sideBar.x+1, sideBar.y+i, '|', nil, tcell.StyleDefault)
	}
	(*manager.screen).Show()
}

func DrawDeleteSideBar(screen tcell.Screen, sideBar SideBar) {
	for i := 0; i < sideBar.size; i++ {
		screen.SetContent(sideBar.x, sideBar.y+i, ' ', nil, tcell.StyleDefault)
		screen.SetContent(sideBar.x+1, sideBar.y+i, ' ', nil, tcell.StyleDefault)
	}
	screen.Show()
}

func (manager *DrawManager) DrawBall(screen tcell.Screen, ball Ball) {
	(*manager.screen).SetContent(ball.x, ball.y, '*', nil, tcell.StyleDefault)
	(*manager.screen).Show()
}

func DrawDeleteBall(screen tcell.Screen, ball Ball) {
	screen.SetContent(ball.x, ball.y, ' ', nil, tcell.StyleDefault)
	screen.Show()
}
