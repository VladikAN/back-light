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

	"github.com/vladikan/back-light/backlight"
)

func main() {
	opt := &backlight.Options{
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

	worker := &backlight.Worker{
		Opt: opt,
		In:  in,
		Ctx: ctx,
	}

	// Find serial port and connect
	go worker.AutoConnect()

	// Capture user screen or draw debug lines
	go func() {
		if opt.IsDebug {
			log.Println("Starting debug mode.\n* Red - top-left.\n* Green - top-right.\n* Blue - bottom-right.\n* Yellow - bottom-left.")
		} else {
			log.Printf("Capturing screen with %d ms refresh rate\n", worker.Opt.RefreshRate)
		}

		tick := time.NewTicker(time.Duration(worker.Opt.RefreshRate * int(time.Millisecond)))
		for {
			select {
			case <-worker.Ctx.Done():
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
				sb := &strings.Builder{}
				for _, c := range rs {
					sb.WriteString(fmt.Sprintf("%02x%02x%02x;", c.R, c.G, c.B))
				}

				worker.In <- sb.String()
			}
		}
	}()

	<-worker.Ctx.Done()
}
