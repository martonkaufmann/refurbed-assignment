package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"refurbed/assignment/internal/notifier"
	"refurbed/assignment/internal/report"
	"syscall"
	"time"

	_ "net/http/pprof"
)

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	cli := flag.NewFlagSet("notify -url=URL [<optional flags>] < messages.txt", flag.ExitOnError)

	flag.Usage = func() {
		fmt.Fprintf(cli.Output(), "Usage %s:\n", os.Args[0])
		cli.PrintDefaults()
	}

	intervalFlag := cli.Uint("interval", 5000, "Interval to check for new messages[ms]")
	timeoutFlag := cli.Uint("timeout", 100, "Request timeout[ms]")
	maxRequestsFlag := cli.Int("requests", 50, "Maximum number of concurrent requests")
	urlFlag := cli.String("url", "", "URL to post messages to")
	cli.Parse(os.Args[1:])

	u, err := url.Parse(*urlFlag)

	if err != nil {
		report.Error(err)
		return
	}

	n := notifier.NewHttpNotifier(
		time.Millisecond*time.Duration(*timeoutFlag),
		int32(*maxRequestsFlag),
		*u,
		ctx,
	)

	go n.Process()

	func() {
		for {
			select {
			case <-shutdown:
				report.Info("Shutdown signal received")
				cancel()
				return
			default:
			}

			fi, err := os.Stdin.Stat()

			if err != nil {
				report.Error(err)
				cancel()
				return
			}

			if fi.Size() > 0 {
				s := bufio.NewScanner(os.Stdin)

				for s.Scan() {
					go n.Enqueue(s.Text())
				}
			}

			time.Sleep(time.Millisecond * time.Duration(*intervalFlag))
		}
	}()

	<-ctx.Done()
}
