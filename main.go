package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/natfaulk/nflogger"
)

const killTime float64 = 30

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

	// open log
	err = openLog()
	if err != nil {
		lg.Fatal("Couldn't open log:", err)
	}
	defer closeLog()

	executable := os.Args[1]
	executableName := filepath.Base(executable)
	executableArgs := os.Args[2:]
	lg.Printf("Executable to run: %s", executableName)

	for i, v := range executableArgs {
		lg.Printf("Arg %d: %s", i, v)
	}

	c := make(chan string)
	lastMessage := time.Now()

	for {
		if isRunning(executableName) {
			lg.Println("IS ALREADY RUNNING")
		} else {
			// check the VPN
			if checkVPN() {
				lg.Println("STARTING EXECUTABLE")

				// receive messages and save last received time
				go func() {
					for msg := range c {
						// else complains about msg being unused
						_ = msg
						deltaT := time.Now().Sub(lastMessage).Seconds()
						lastMessage = time.Now()
						lg.Println(deltaT)
					}
				}()

				// kill program if it is not responding
				go func() {
					for true {
						if isRunning(executableName) {
							deltaT := time.Now().Sub(lastMessage).Seconds()
							if deltaT > killTime {
								lg.Println("THIS RAN")
								err := runWithOutput(nil, "pkill", executableName)
								if err != nil {
									lg.Println("Error running pkill")
									lg.Println(err)
								}
							}
						}

						time.Sleep(time.Second * 60)
					}
				}()

				// reset the last message count
				lastMessage = time.Now()
				// lets start the executable
				err := runWithOutput(c, executable, executableArgs...)
				if err != nil {
					lg.Printf("Error running executable %s", executable)
					lg.Println(err)
				}
			} else {
				lg.Println("VPN NOT CONNECTED")
				if err := restartVPN(); err != nil {
					lg.Println("Failed to restart VPN")
					lg.Println(err)
				}
			}
		}

		// sleep for a minute regardless what happened so we don't go crazy
		// if something fails very quickly...
		lg.Println("SLEEPING 1 MIN")
		time.Sleep(1 * time.Minute)
	}
}
