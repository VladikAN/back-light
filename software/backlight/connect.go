package backlight

import (
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
)

func (worker *Worker) AutoConnect() {
	for {
		select {
		case <-worker.Ctx.Done():
			log.Printf("Serial channel was closed\n")
			return
		default:
			err := connect(worker.In)
			if err != nil {
				log.Printf("ERROR unable to send data: %s\n", err)
				log.Printf("Reconnect in %d ms\n", worker.Opt.Timeout)
				time.Sleep(time.Duration(worker.Opt.Timeout * time.Hour.Milliseconds()))
			}
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
		// log.Printf("- %s\n", val)
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
