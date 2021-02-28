package main

import "os"

func checkEnvVarExists(name string) bool {
	localAddr := os.Getenv(name)
	if localAddr == "" {
		lg.Printf("Failed to get environment variable %s", name)
		lg.Println("Have you added a .env file to the executable directory?")

		return false
	}
	return true
}
