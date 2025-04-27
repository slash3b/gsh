package main

import (
    "context"
    "fmt"
    "bufio"
    "os"
    "strings"
    // "io/fs"
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

    _, _ = ctx, cancel

    // os.StartProcess use it to start users programs


    sc := bufio.NewScanner(os.Stdin)

    printPromptStart()
    for sc.Scan() {
        processInput(sc.Text())

        printPromptStart()
    }
}

func printPromptStart() {
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
