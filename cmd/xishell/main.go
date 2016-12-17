package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/alexsward/xisdb"
	"github.com/alexsward/xisdb/ql"
)

func main() {
	interrupt()

	out := os.Stdout
	io.WriteString(out, "xisdb shell, default will be port 7632 or something\n")

	xis, err := xisdb.Open(&xisdb.Options{InMemory: true})
	if err != nil {
		fmt.Printf("Error opening database: %s\n", err)
		return
	}

	qe := &xisdb.QueryEngine{}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		io.WriteString(out, "> ")
		if !scanner.Scan() {
			return
		}

		input := scanner.Text()
		if handleShellCommands(input) {
			continue
		}

		statements, err := ql.Parse(input)
		if err != nil {
			io.WriteString(out, fmt.Sprintf("Error parsing statement: %s\n", err))
		}
		ctx := newContext(xis)
		err = qe.Execute(statements, ctx)
		if err != nil {
			io.WriteString(out, fmt.Sprintf("Error executing statements: %s\n", err))
		}
		t := time.NewTicker(time.Millisecond * 100)
		select {
		case <-t.C:
			io.WriteString(out, "Timed out\n")
			return
		case r := <-ctx.Results:
			io.WriteString(out, fmt.Sprintf("Received:[%s]\n", r))
		}
		t.Stop()
	}
}

func handleShellCommands(input string) bool {
	switch strings.ToLower(input) {
	case "quit", "exit":
		io.WriteString(os.Stdout, "Quitting xisdb shell\n")
		os.Exit(1)
		return true
	default:
		return false
	}
}

func newContext(db *xisdb.DB) *xisdb.QueryEngineContext {
	return &xisdb.QueryEngineContext{
		DB:      db,
		Results: make(chan xisdb.Item, 0),
	}
}

func interrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()
}
