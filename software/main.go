package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.bug.st/serial"
)

func main() {
	input := make(chan string)

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		close(input)
	}()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		autoConnect(input)
	}()

	input <- "Test 1"
	input <- "Test 2"
	input <- "Test 3"

	wg.Wait()
}

func autoConnect(input <-chan string) {
	for {
		err := connect(input)
		if err != nil {
			log.Printf("ERROR unable to send data: %s\n", err)
			log.Printf("Reconnect in 5 seconds")
			time.Sleep(5 * time.Second)
		} else {
			log.Printf("Serial channel was closed\n")
			break
		}
	}
}

func connect(input <-chan string) error {
	port, err := getPort()
	if err != nil {
		return err
	}

	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	srl, err := serial.Open(port, mode)
	if err != nil {
		return err
	}
	defer srl.Close()

	log.Printf("Connected to %s port\n", port)
	for val := range input {
		_, err := srl.Write([]byte(val))
		if err != nil {
			return err
		}
	}

	return nil
}

func getPort() (string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return "", err
	}

	if ln := len(ports); ln == 0 {
		return "", fmt.Errorf("no serial ports available")
	} else if ln > 1 {
		log.Printf("WARN More than 1 serial port found. First in the list will be used\n")
	}

	return ports[0], nil
}
