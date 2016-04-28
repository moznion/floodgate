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
	Help      bool   `cli:"h,help" usage:"display help"`
	Interval  int    `cli:"i,interval" usage:"intarval time to flush (second)" dft:"0"`
	Threshold int    `cli:"t,threshold" usage:"throshold size of memory to flush (byte)" dft:"0"`
	Concat    string `cli:"c,concat" usage:"character to concat for each line" dft:"\n"`
	IsStderr  bool   `cli:"stderr" usage:"flush to STDERR"`
}

type floodgate struct {
	concat []byte
	buf    [][]byte
	size   int
	dst    *os.File
	mut    *sync.Mutex
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

		fg := &floodgate{
			concat: []byte(argv.Concat),
			mut:    new(sync.Mutex),
			size:   0,
			dst:    os.Stdout,
		}
		if argv.IsStderr {
			fg.dst = os.Stderr
		}

		var flusher func(int)
		if argv.Threshold > 0 {
			flusher = fg.flush
		}
		go fg.scan(argv.Threshold, flusher)

		if argv.Interval > 0 {
			go fg.tick(argv.Interval)
		}

		return nil
	})

	select {}
}

func (fg *floodgate) tick(interval int) {
	for _ = range time.Tick(time.Duration(interval) * time.Second) {
		fg.flush(0)
	}
}

func (fg *floodgate) scan(threshold int, flusher func(int)) {
	r := bufio.NewReader(os.Stdin)
	concatStr := string(fg.concat)

	for {
		lineBytes, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				os.Exit(0)
			}
			panic(err)
		}

		fg.size += len(lineBytes)
		if concatStr != "\n" {
			lineBytes = append(lineBytes[:len(lineBytes)-1], fg.concat...)
		}
		fg.buf = append(fg.buf, lineBytes)

		if flusher != nil {
			flusher(threshold)
		}
	}
}

func (fg *floodgate) flush(tsize int) {
	fg.mut.Lock()
	if fg.size > tsize {
		for _, buf := range fg.buf {
			fmt.Fprintf(fg.dst, string(buf))
		}
		fg.buf = fg.buf[:0]
		fg.size = 0
	}
	fg.mut.Unlock()
}
