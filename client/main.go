package main

import (
	"flag"
	"log"

	"github.com/gdamore/tcell/v2"
)

var LOG_Y_POS int = 50

func main() {
	Flag := flag.Int("LOG", 0, "use log")
	flag.Parse()
	if *Flag == 1 {
		onLog = true
	}

	quit := make(chan struct{})

	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	if err = screen.Init(); err != nil {
		log.Fatal(err)
	}
	defer screen.Fini()

	width := 160
	height := 45

	InitDrawManager(&screen)

	DebugLog("debug")

	// command loop
	commandChan := make(chan tcell.Key, 10)
	commandReady := make(chan int, 1)
	go func() {
		for {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape {
					close(quit)
				} else if ev.Key() == tcell.KeyUp {
					commandChan <- ev.Key()
				} else if ev.Key() == tcell.KeyDown {
					commandChan <- ev.Key()
				}
				<-commandReady
			}
		}
	}()

	// マッチング処理
	tcpconn, err := NewConnection()
	if err != nil {
		DrawError(screen, err.Error())
	}

	go tcpconn.Process()

	InitGameManager(tcpconn, commandChan, commandReady)

	if err := gameManager.matchingStart(); err != nil {
		return
	}

	// 画面表示
	DrawFrame(screen, height, width)

	go gameManager.process()

	<-quit
}
