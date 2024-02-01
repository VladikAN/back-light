package backlight

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
)

func (w *Worker) Connect(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			w.Abort()
			return
		default:
			err := w.connectSerial()
			if err != nil {
				log.Printf("ERROR unable to send data: %s\n", err)
				log.Printf("Reconnect in %d ms\n", w.Opt.ConnectTimeout)
				time.Sleep(time.Duration(w.Opt.ConnectTimeout * time.Hour.Milliseconds()))
			}
		}
	}
}

func (w *Worker) Abort() {
	if w.IsReady {
		w.In <- w.ToSerial(w.DrawEmpty())
	}

	close(w.In)
}

func (w *Worker) connectSerial() error {
	port, err := getPort()
	if err != nil {
		return err
	}

	mode := &serial.Mode{
		BaudRate: w.Opt.SerialSpeed,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	srl, err := serial.Open(port, mode)
	if err != nil {
		return err
	}
	defer srl.Close()

	w.IsReady = true
	defer func() {
		w.IsReady = false
	}()

	log.Printf("Connected to %s port\n", port)
	for val := range w.In {
		// log.Printf("- %s\n", val)
		_, err := srl.Write([]byte(val))
		if err != nil {
			return err
		}
	}

	// Let serial complete the transfer
	time.Sleep(time.Duration(100 * time.Hour.Milliseconds()))
	log.Printf("Serial channel was closed\n")
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
		log.Printf("WARN %d serial ports found, %s port is used\n", len(ports), ports[0])
	}

	return ports[0], nil
}
