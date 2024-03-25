package main

type SideBar struct {
	x      int // 最上部のx座標
	y      int // 最上部のy座標
	size   int // バーのサイズ
	top    int
	bottom int
}

func NewSideBar(x int, y int, size int, heigt int) *SideBar {
	top := 0 + 3
	bottom := heigt - 1 - size
	sideBar := SideBar{x, y, size, top, bottom}
	return &sideBar
}

func (bar *SideBar) SlideUpSideBar(y int) {
	if y < 0 {
		return
	}

	moved := bar.y - y
	if moved < bar.top {
		moved = bar.top
	}
	bar.y = moved
}

func (bar *SideBar) SlideDownSideBar(y int) {
	if y < 0 {
		return
	}

	moved := bar.y + y
	if moved > bar.bottom {
		moved = bar.bottom
	}
	bar.y = moved
}
