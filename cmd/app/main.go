package main

import "indicators-monitor/internal/watcher"

func main() {
	block := make(chan bool)
	watcher.NewWatcher([]string{
		"ETH", "BTC", "DOGE", "SOL",
		"BNB", "ADA", "SAND", "MATIC",
		"OP",
	}).Start()

	<-block
}
