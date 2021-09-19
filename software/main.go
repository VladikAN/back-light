package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.bug.st/serial"
)

type Options struct {
	// LED line width, how many LED's used for top and bottom parts of the screen. Can be formed in scuare or circle.
	Width int
	// LED line height, how many LED's used for left and right parts of the screen. Can be formed in scuare or circle.
	Height int
	// Debug mode will use constant colors, ie red for top and bottom, green for left and right, yellow for the corners.
	IsDebug bool
	// Time ms between screen captures and data transfer to serial port.
	RefreshRate int
	// Time ms between serial port searches
	Timeout int64
}

type worker struct {
	opt *Options
	in  chan string
	ctx context.Context
}

func main() {
	opt := &Options{
		Width:       4,
		Height:      4,
		IsDebug:     true,
		RefreshRate: 100,
		Timeout:     1000,
	}

	in := make(chan string)

	// Handle terminate signal like CTRL+C
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		cancel()
	}()

	worker := &worker{
		opt: opt,
		in:  in,
		ctx: ctx,
	}
	wg := sync.WaitGroup{}

	// Find serial port and connect
	wg.Add(1)
	go func() {
		defer wg.Done()
		worker.autoConnect()
	}()

	// Capture user screen or draw debug lines
	wg.Add(1)
	go func() {
		defer wg.Done()

		if opt.IsDebug {
			log.Println("Starting debug mode.\n* Yellow - corners.\n* Red - top and bottom.\n* Green - left and right.")
		} else {
			log.Printf("Capturing screen with %d ms refresh rate\n", worker.opt.RefreshRate)
		}

		tick := time.NewTicker(time.Duration(worker.opt.RefreshRate * int(time.Millisecond)))
		for {
			select {
			case <-worker.ctx.Done():
				close(worker.in)
				return // exit if ctx is done
			case <-tick.C:
				if opt.IsDebug {
					worker.drawDebug()
				} else {
					worker.drawScreen()
				}
			}
		}
	}()

	wg.Wait()
}

func (worker *worker) drawDebug() {
	log.Println("TICK")
}

func (worker *worker) drawScreen() {

}

func (worker *worker) autoConnect() {
	for {
		select {
		case <-worker.ctx.Done():
			log.Printf("Serial channel was closed\n")
			return
		default:
			err := connect(worker.in)
			if err != nil {
				log.Printf("ERROR unable to send data: %s\n", err)
				log.Printf("Reconnect in %d ms\n", worker.opt.Timeout)
				time.Sleep(time.Duration(worker.opt.Timeout * time.Hour.Milliseconds()))
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
