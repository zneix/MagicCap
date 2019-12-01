package core

import (
	"fmt"
	"github.com/hackebrot/turtle"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/gobuffalo/packr"
)

var (
	// HomeDir defines the home directory.
	HomeDir, _ = os.UserHomeDir()

	// ConfigPath defines the MagicCap folder path.
	ConfigPath = path.Join(HomeDir, ".magiccap")

	// Assets contains all of the data from "assets" when compiled.
	Assets = packr.NewBox("../assets")

	// Version defines the version.
	Version = "3.0.0a1"

	// Emojis are all usable emojis.
	Emojis = make([]string, 0)
)

// Start is the main entrypoint for the application.
func Start() {
	// Handle the random seed.
	rand.Seed(time.Now().UnixNano())

	// Load in all of the emojis.
	for  _, value := range turtle.Emojis {
		Emojis = append(Emojis, value.String())
	}

	// Gets the start time.
	StartTime := time.Now()

	// Ensures that ConfigPath exists.
	_ = os.MkdirAll(ConfigPath, 0777)

	// Does any migrations which are needed.
	MigrateFrom2()

	// Boot message.
	println("MagicCap " + Version + " - Copyright (C) MagicCap Development Team 2018-2019.")

	// Loads up the uploader kernel.
	LoadUploadersKernel()

	// Loads the SQLite3 DB.
	LoadDatabase()

	// Starts the tray.
	RestartTrayProcess()

	// Defines how long it took.
	elapsed := time.Since(StartTime)
	fmt.Printf("Initialisation took %s.\n", elapsed)

	// Keep the app alive.
	for {
		time.Sleep(time.Second)
	}
	// TODO: Make a install ID.
}
