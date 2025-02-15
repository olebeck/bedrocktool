package utils

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bedrock-tool/bedrocktool/locale"
	"github.com/google/subcommands"
	"github.com/sirupsen/logrus"
)

type MessageResponse struct {
	Ok   bool
	Data interface{}
}

type UI interface {
	Init() bool
	Start(context.Context, context.CancelFunc) error
	Message(name string, data interface{}) MessageResponse
	ServerInput(context.Context, string) (string, string, error)
}

type BaseUI struct {
	UI
}

func (u *BaseUI) Message(name string, data interface{}) MessageResponse {
	return MessageResponse{
		Ok:   false,
		Data: nil,
	}
}

func (u *BaseUI) ServerInput(ctx context.Context, server string) (string, string, error) {
	address, name, err := ServerInput(ctx, server)
	return address, name, err
}

var currentUI UI

func SetCurrentUI(ui UI) {
	currentUI = ui
}

type InteractiveCLI struct {
	BaseUI
}

func (c *InteractiveCLI) Init() bool {
	currentUI = c
	return true
}

func (c *InteractiveCLI) Start(ctx context.Context, cancel context.CancelFunc) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		fmt.Println(locale.Loc("available_commands", nil))
		for name, cmd := range ValidCMDs {
			fmt.Printf("\t%s\t%s\n", name, cmd.Synopsis())
		}
		fmt.Println(locale.Loc("use_to_run_command", nil))

		cmd, cancelled := UserInput(ctx, locale.Loc("input_command", nil))
		if cancelled {
			cancel()
			return nil
		}
		_cmd := strings.Split(cmd, " ")
		os.Args = append(os.Args, _cmd...)
	}
	flag.Parse()

	InitDNS()
	InitExtraDebug()

	subcommands.Execute(ctx)

	if Options.IsInteractive {
		logrus.Info(locale.Loc("enter_to_exit", nil))
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
	}
	return nil
}

var MakeGui = func() UI {
	return &InteractiveCLI{}
}
