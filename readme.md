####
gsh is a toy shell that supports basic functionality
todo: add list and examples maybe

```
/tmp/gsh ⮕ ls
.git/ .gitignore LICENSE go.mod gsh main.go readme.md
/tmp/gsh ⮕ cat go.mod
module github.com/slash3b/gsh

go 1.24.2

/home/slash3b/Projects/personal/gsh ⮕
```

#### How to install

go install github.com/slash3b/gsh@latest

#### Misc interesting stuff:

##### How Ctrl+D Works in Terminals
The terminal driver operates in cooked mode by default, buffering input until a newline or EOF is received.
When Ctrl+D is pressed at the start of a line, the terminal driver immediately returns the buffered input (which is empty) to the program, signaling EOF.
If Ctrl+D is pressed after typing some characters (not at the start of a line), it causes the terminal to send the current buffered input to the program without a newline, not EOF. This means you might need to press Ctrl+D twice to signal EOF if you have typed partial input
