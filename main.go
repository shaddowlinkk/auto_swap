package main

import (
	"bufio"
	"context"
	"fmt"
	"hotReloadplugins/Utils"
	"os"
	"strings"
	"time"
)

type EnvConfig struct {
	Watch       string `json:"Watch"`
	Dest        string `json:"dest"`
	StartServer struct {
		Type    string   `json:"type"`
		Command string   `json:"command"`
		Args    []string `json:"args"`
	} `json:"start_server"`
	StopServer struct {
		Type    string   `json:"type"`
		Command string   `json:"command"`
		Args    []string `json:"args"`
	} `json:"stop_server"`
}

func ReadServerConfig(filename string) EnvConfig {
	return Utils.ReadConfig[EnvConfig](filename)
}

func main() {
	config := ReadServerConfig("config.json")
	ctx, can := context.WithCancel(context.Background())
	term := Term{Ctx: ctx, Cancel: can}

	main_ctx, main_can := context.WithCancel(context.Background())

	go func() {
		errs := RunCommand(config.StartServer.Type, config.StartServer.Command, config.StartServer.Args, &term)
		if errs != nil {
			can()
			return
		}
		go func() {
			for {

				reader := bufio.NewReader(os.Stdin)

				// Read the input string until a newline character
				inputStr, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Error reading input:", err)
					return
				}
				if strings.Contains(inputStr, "ctxstop") {
					main_can()
				}
				err = RunCommand("STDIN", inputStr, nil, &term)
				if err != nil {
					return
				}
			}

		}()
		var old int64 = -1
		for {
			i, err := os.Stat(config.Watch)
			if err == nil && i.Size() != old {
				old = i.Size()
				fmt.Printf("old:%d,new%d\n", old, i.Size())
				time.Sleep(1 * time.Second)
				/*			_, err := io.WriteString(*term.SendCommand, "echo \"fart\"\r\n")
							if err != nil {
								return
							}*/
				err = RunCommand(config.StopServer.Type, config.StopServer.Command, config.StopServer.Args, &term)
				if err != nil {
					return
				}
				done := false
				for !done {
					select {
					case <-term.Ctx.Done():
						ctx, can := context.WithCancel(context.Background())
						term = Term{Ctx: ctx, Cancel: can}
						errs := RunCommand(config.StartServer.Type, config.StartServer.Command, config.StartServer.Args, &term)
						if errs != nil {
							can()
							return
						}
						done = true
						break
					}
				}

				fmt.Println("loading new file")
				err = os.Rename(config.Watch, config.Dest)
				if err != nil {
					fmt.Println(err)
				}
				err = nil
				time.Sleep(2 * time.Second)
				fmt.Println("Done Restart")

			} else if os.IsNotExist(err) {
				continue
			} else {
				/*				can()
								fmt.Println("Error:", err)*/
			}
		}
	}()
	for {
		select {
		case <-main_ctx.Done():
			can()
			fmt.Println("Stoping")
			return
		}
	}
}
