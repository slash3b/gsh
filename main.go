package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	// run this program
	// support a few builtins ls, cd
	// ctrl-c should terminate running program by shell
	// EOF or ctrl D or exit should terminate shell itself
	// implement repl loop
	//  take string after enter, how to do it??
	//  execute it

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigch := make(chan os.Signal)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

	// cancelling goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case sg := <-sigch:
				fmt.Println("caught signal", sg)
				cancel()
			}
		}
	}()

	sc := bufio.NewScanner(os.Stdin)

	// input goroutine
	inputch := make(chan string)
	go func() {
		// todo(ilya): should be a better way to "close" this?
		// check source code of Scan
		for sc.Scan() {
			inputch <- sc.Text()
		}
	}()

	printPromptStart()
	for {
		select {
		case <-ctx.Done():
			return
		case in := <-inputch:
			processInput(in)
			printPromptStart()
		}
	}
}

func printPromptStart() {
	// show pwd as well
	fmt.Print(">> ")
}

func processInput(s string) {
	s = strings.TrimSpace(s)
	switch s {
	case "":
		return
	case "ls":
		listFiles()
	case "exit":
		os.Exit(0)
	default:
		// os.StartProcess use it to start users programs
		// this will actually run
		fmt.Println(s)
	}
}

func listFiles() {
	entries, err := os.ReadDir(".")
	if err != nil {
		fmt.Println(err)

		return
	}

	for _, v := range entries {
		if v.IsDir() {
			fmt.Print(v.Name()+"/", " ")

			continue
		}

		fmt.Print(v.Name(), " ")
	}

	fmt.Println()
}
