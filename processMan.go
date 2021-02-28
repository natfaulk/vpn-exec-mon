package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

func isRunning(name string) bool {
	cmd := exec.Command("pgrep", name)
	var outb bytes.Buffer
	cmd.Stdout = &outb

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				return false
			}
		}

		fmt.Println("An error ocurred")
		log.Fatal(err)
	}

	// fmt.Printf(outb.String())

	return true
}

func runWithOutput(name string, arg ...string) {
	cmd := exec.Command(name, arg...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("An error ocurred")
		log.Fatal(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("An error ocurred")
		log.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("An error ocurred")
		log.Fatal(err)
	}

	merged := io.MultiReader(stderr, stdout)
	scanner := bufio.NewScanner(merged)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println("An error ocurred")
		log.Fatal(err)
	}
}

// canPing returns if pinging the target was successful
// if not returns the status code - sometimes we want a ping to not work...
// if returns true the status code will be 0
// shouldnt return false without a status code, but if it couldnt get the status code will return -1
func canPing(address string) (bool, int) {
	var nAttempts int = 1
	var timeoutSecs int = 5

	cmd := exec.Command("ping", address, "-c", fmt.Sprint(nAttempts), "-W", fmt.Sprint(timeoutSecs))
	var outb bytes.Buffer
	cmd.Stdout = &outb

	err := cmd.Run()
	if err != nil {
		// check it wasn't another error (like no connectivity)
		if exitError, ok := err.(*exec.ExitError); ok {
			return false, exitError.ExitCode()
		}

		// Hopefully shouldnt run but you never know...
		return false, -1
	}

	return true, 0
}

func checkVPN() bool {
	//  is there general connectivity
	ok, _ := canPing(os.Getenv("REMOTE_ADDR"))
	if !ok {
		return false
	}

	// counterintuitively, if this ping was successful it means we could
	// ping the router (which we shouldn't be able to) therefore the vpn is not up
	ok, statusCode := canPing(os.Getenv("LOCAL_ADDR"))
	if ok {
		return false
	}

	// if there is any error other than return 1 something probs went wrong...
	if statusCode != 1 {
		return false
	}

	return true
}

func stopVPN() {
	runWithOutput("./scripts/vpnend.sh")
}

func startVPN() {
	runWithOutput("./scripts/vpnstart.sh")
}

func restartVPN() {
	lg.Println("Stopping VPN...")
	stopVPN()
	time.Sleep(time.Second * 5)
	lg.Println("Starting VPN...")
	startVPN()
}
