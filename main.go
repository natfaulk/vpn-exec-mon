package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/natfaulk/nflogger"
)

var lg *log.Logger = nflogger.Make("main")

func main() {
	err := godotenv.Load()
	if err != nil {
		lg.Print("Error loading .env file")
	}

	if len(os.Args) < 2 {
		lg.Println("At least one arguement is required (the executable to run).")
		lg.Println("Add the arguements for the executable to run as subsequent arguements")
		lg.Fatal("e.g.: ./vpn-exec-mon myprogram arg1 arg2")
	}

	if !checkEnvVarExists("LOCAL_ADDR") {
		os.Exit(1)
	}

	if !checkEnvVarExists("REMOTE_ADDR") {
		os.Exit(1)
	}

	executable := os.Args[1]
	executableName := filepath.Base(executable)
	executableArgs := os.Args[2:]
	lg.Printf("Executable to run: %s", executableName)

	for i, v := range executableArgs {
		lg.Printf("Arg %d: %s", i, v)
	}

	for {
		if isRunning(executableName) {
			lg.Println("IS ALREADY RUNNING")
		} else {
			// check the VPN
			if checkVPN() {
				lg.Println("STARTING ExECUTABLE")
				// lets start the executable
				runWithOutput(executable, executableArgs...)
			} else {
				lg.Println("VPN NOT CONNECTED")
				restartVPN()
			}
		}

		// sleep for a minute regardless what happened so we don't go crazy
		// if something fails very quickly...
		lg.Println("SLEEPING 1 MIN")
		time.Sleep(1 * time.Minute)
	}
}
