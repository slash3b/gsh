package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
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

	// todo: implement builtin s, like exit, cd and etc.

	// todo: handle case like gsh --version or gsh version

	// todo: how the heck clear works out of the box???

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigch := make(chan os.Signal)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGHUP)

	// just a courtesy
	// we could have just abandoned
	var wg sync.WaitGroup
	wg.Add(2)

	// cancelling goroutine
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-sigch:
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
			select {
			case inputch <- sc.Text():
			case <-ctx.Done():
				break
			}
		}
		// Ctrl+D will be interpreted as io.EOF marker
		// and sc.Err will be nil
		if sc.Err() != nil {
			fmt.Println("error:", sc.Err())
		}

		cancel()
	}()

	printPromptStart()

	for {
		select {
		case <-ctx.Done():
			return
		case in := <-inputch:
			processInput(ctx, in)
			printPromptStart()
		}
	}

	wg.Wait()
}

func printPromptStart() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// show pwd as well
	fmt.Print(wd + "\033[1m 🐚 \033[0m")
}

func processInput(ctx context.Context, s string) {
	s = strings.TrimSpace(s)

	if s == "" {
		return
	}

	// is builtin

	switch s {
	case "":
		return
	case "ls":
		listFiles()
	case "exit":
		os.Exit(0)
	default:
		// tokenize

		tokens := strings.Fields(s)

		cmd := exec.CommandContext(ctx, tokens[0], tokens[1:]...)

		b, err := cmd.Output()
		if err != nil {
			fmt.Printf("output error: %s\n", err)
			return
		}

		fmt.Println(string(b))
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
