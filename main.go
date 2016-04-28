package floodgate

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/mkideal/cli"
)

type opt struct {
	Help      bool `cli:"h,help" usage:"display help"`
	Interval  int  `cli:"i,interval" usage:"intarval time to flush (second)" dft:"0"`
	Threshold int  `cli:"t,threshold" usage:"throshold size of memory to flush (byte)" dft:"0"`
}

type buffer struct {
	buf  [][]byte
	size int
	mut  *sync.Mutex
}

// Run floodgate application
func Run(args []string) {
	cli.Run(&opt{}, func(ctx *cli.Context) error {
		argv := ctx.Argv().(*opt)
		if argv.Help {
			ctx.String(ctx.Usage())
			os.Exit(0)
		}

		if argv.Interval <= 0 && argv.Threshold <= 0 {
			fmt.Fprintf(os.Stderr, "[ERROR] Interval or threshold must be at least 0\n")
			ctx.String(ctx.Usage())
			os.Exit(2)
		}

		b := &buffer{
			mut:  new(sync.Mutex),
			size: 0,
		}

		var flusher func(int)
		if argv.Threshold > 0 {
			flusher = b.flushByThreshold
		}
		go b.scan(argv.Threshold, flusher)

		if argv.Interval > 0 {
			go b.tick(argv.Interval)
		}

		return nil
	})

	select {}
}

func (b *buffer) tick(interval int) {
	for _ = range time.Tick(time.Duration(interval) * time.Second) {
		b.mut.Lock()
		if b.size > 0 {
			for _, buf := range b.buf {
				fmt.Print(string(buf))
			}
			b.buf = b.buf[:0]
		}
		b.mut.Unlock()
	}
}

func (b *buffer) scan(threshold int, flusher func(int)) {
	r := bufio.NewReader(os.Stdin)

	for {
		lineBytes, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				os.Exit(0)
			}
			panic(err)
		}

		b.size += len(lineBytes)
		b.buf = append(b.buf, lineBytes)

		if flusher != nil {
			flusher(threshold)
		}
	}
}

func (b *buffer) flushByThreshold(threshold int) {
	b.mut.Lock()
	if b.size >= threshold {
		for _, buf := range b.buf {
			fmt.Print(string(buf))
		}
		b.buf = b.buf[:0]
	}
	b.mut.Unlock()
}
