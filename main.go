package main

import (
	"bufio"
	"os"

	"strings"

	"flag"

	"github.com/kpango/glg"
	"github.com/vetcher/easystarter/backend"
	"github.com/vetcher/easystarter/commands"
)

// TODO: specify service version
// TODO: add cleaning command
// TODO: open logs

const (
	VERSION          = "0.2"
	WelcomeTip       = "Easy Starter " + VERSION
	MKDIR_PERMISSION = 0777

	CMD_START   = "start"
	CMD_STOP    = "stop"
	CMD_RESTART = "restart"
	CMD_PS      = "ps"
	CMD_ENV     = "env"
	CMD_EXIT    = "exit"
	CMD_VERSION = "version"
	CMD_KILL    = "kill"

	EXIT_CODE_SETUP_ENV_ERR = 1 + iota
	EXIT_CODE_INIT_LOGS_DIR_ERR
)

func init() {
	if !backend.SetupEnv() {
		glg.Fatal("I'm out, can't setup env")
		os.Exit(EXIT_CODE_SETUP_ENV_ERR)
	}
	_, err := os.Stat("logs")
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("logs", MKDIR_PERMISSION)
		} else {
			glg.Fatal(err)
			os.Exit(EXIT_CODE_INIT_LOGS_DIR_ERR)
		}
	}
}

func main() {
	allCommands := map[string]commands.Command{
		CMD_START:   &commands.StartCommand{},
		CMD_STOP:    &commands.StopCommand{},
		CMD_PS:      &commands.PSCommand{},
		CMD_ENV:     &commands.EnvCommand{},
		CMD_RESTART: &commands.RestartCommand{},
		CMD_VERSION: &commands.VersionCommand{VERSION},
		CMD_EXIT:    &commands.ExitCommand{},
		"":          &commands.EmptyCommand{},
		CMD_KILL:    &commands.KillCommand{},
	}
	flag.Parse()
	glg.Print(WelcomeTip)
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		text := stdin.Text()
		inputCommands := strings.Split(text, " ")
		command, ok := allCommands[inputCommands[0]]
		if ok {
			err := command.Validate(inputCommands[1:]...)
			if err != nil {
				glg.Errorf("Validation error: %v", err)
				continue
			}
			err = command.Exec(inputCommands[1:]...)
			if err != nil {
				glg.Error(err)
				return
			}
		} else {
			glg.Printf("`%v` is wrong command, try to `help`.", inputCommands[0])
		}
	}
}
