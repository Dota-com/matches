package main

import "matches/internal/config"

func main() {
	conf := config.MustLoad()
	log := config.SetupLoger(conf.Env)
	log.Debug("Debug running")
}
