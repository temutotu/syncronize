package main

import (
	"log"
	"os"
)

const FILE string = "clientLog.txt"

var onLog bool = false

func DebugLog(data any) {
	if !onLog {
		return
	}
	file, err := os.OpenFile(FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// エラーハンドリング
		log.Fatal(err)
	}

	// 関数が終了する際にファイルを閉じる
	defer file.Close()

	// ログの出力先をファイルに設定
	log.SetOutput(file)

	// ログを出力
	log.Println(data)
}
