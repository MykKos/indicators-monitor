package main

import (
	"indicators-monitor/internal/watcher"
	"os"
	"strings"
)

func main() {
	block := make(chan bool)
	// watcher.NewWatcher([]string{
	// 	"ETH", "BTC", "DOGE", "SOL",
	// 	"BNB", "ADA", "SAND", "MATIC",
	// 	"OP",
	// }).Start()

	watcher.NewWatcher(readTokens())

	<-block
}

func readTokens() []string {
	f, _ := os.ReadFile("configs/tokens-list")
	return strings.Split(string(f), "\n")
}
