package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Term struct {
	SendCommand *io.WriteCloser
	Ctx         context.Context
	Cancel      context.CancelFunc
}

func RunCommand(Type string, command string, Args []string, cmd *Term) error {
	switch Type {
	case "STDIN":
		//fmt.Printf("%s%s\n", command, strings.Join(Args, " "))
		if len(Args) > 0 {
		} else {
			_, err := io.WriteString(*cmd.SendCommand, fmt.Sprintf("%s\n", command))
			if err != nil {
				return err
			}
		}

	case "CMD":
		go func() {

			ncmd := exec.CommandContext(cmd.Ctx, "cmd.exe", "/C", command)
			stdin, er := ncmd.StdinPipe()
			ncmd.Stdout = os.Stdout
			if er != nil {
				fmt.Println(er)
				panic(" failed get pipe")
			}
			cmd.SendCommand = &stdin
			var err error
			err = ncmd.Run()
			if err != nil {
				fmt.Printf("Run Error:%s\n", err)
				panic(" failed top run command")
			}
			fmt.Println("done")
			cmd.Cancel()
		}()
	default:
		return fmt.Errorf("error Type not found")

	}
	return nil
}
