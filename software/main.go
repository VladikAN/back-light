package main

import (
	"context"
	"image/color"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/vladikan/back-light/backlight"
)

func main() {
	opt := &backlight.Options{
		Width:          3,
		Height:         3,
		IsDebug:        false,
		RefreshRate:    33,
		SerialSpeed:    9600,
		ConnectTimeout: 1000,
		ColorLimit:     32,
		Invert:         false,
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

	worker := &backlight.Worker{
		Opt: opt,
		In:  in,
		Ctx: ctx,
	}

	wg := &sync.WaitGroup{}

	// Find serial port and connect
	wg.Add(1)
	go func() {
		defer wg.Done()
		worker.AutoConnect()
	}()

	// Capture user screen or draw debug lines
	wg.Add(1)
	go func() {
		defer wg.Done()
		if opt.IsDebug {
			log.Println("Starting debug mode.\n* Red - top-left.\n* Green - top-right.\n* Blue - bottom-right.\n* Yellow - bottom-left.")
		} else {
			log.Printf("Capturing screen with %d ms refresh rate\n", worker.Opt.RefreshRate)
		}

		tick := time.NewTicker(time.Duration(worker.Opt.RefreshRate * int(time.Millisecond)))
		for {
			select {
			case <-worker.Ctx.Done():
				worker.In <- worker.ToSerial(worker.DrawEmpty())
				close(worker.In)
				return // exit if ctx is done
			case <-tick.C:
				var rs []color.RGBA
				var err error
				if opt.IsDebug {
					rs = worker.DrawDebug()
				} else {
					rs, err = worker.DrawScreen()
				}

				if err != nil {
					log.Printf("Error occurred %s\n", err)
					continue
				}

				rs = worker.FilterOutput(rs)
				worker.In <- worker.ToSerial(rs)
			}
		}
	}()

	wg.Wait()
}
