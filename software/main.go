package main

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.bug.st/serial"
)

type Options struct {
	// LED line width, how many LED's used for top and bottom parts of the screen. Can be formed in square or circle.
	Width int
	// LED line height, how many LED's used for left and right parts of the screen. Can be formed in square or circle.
	Height int
	// Draw debug with constant colors. Red - top-left, Green - top-right, Blue - bottom-right and Yellow - bottom-left.
	IsDebug bool
	// Time ms between screen captures and data transfer to serial port.
	RefreshRate int
	// Time ms between serial port searches
	Timeout int64
	// Invert axis
	Invert bool
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
		IsDebug:     false,
		RefreshRate: 100,
		Timeout:     1000,
		Invert:      true,
	}

	in := make(chan string, 1)

	// Handle terminate signal, ex CTRL+C
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

	// Find serial port and connect
	go worker.autoConnect()

	// Capture user screen or draw debug lines
	go func() {
		if opt.IsDebug {
			log.Println("Starting debug mode.\n* Red - top-left.\n* Green - top-right.\n* Blue - bottom-right.\n* Yellow - bottom-left.")
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
				var rs []color.RGBA
				var err error
				if opt.IsDebug {
					rs = worker.drawDebug()
				} else {
					rs, err = worker.drawScreen()
				}

				if err != nil {
					log.Printf("Error occurred %s\n", err)
					continue
				}

				rs = worker.filterOutput(rs)
				sb := &strings.Builder{}
				for _, c := range rs {
					sb.WriteString(fmt.Sprintf("%02x%02x%02x;", c.R, c.G, c.B))
				}

				worker.in <- sb.String()
			}
		}
	}()

	<-worker.ctx.Done()
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
