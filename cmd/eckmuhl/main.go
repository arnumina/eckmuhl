/*
#######
##                __               __   __
##       ___ ____/ /__ __ _  __ __/ /  / /
##      / -_) __/  '_//  ' \/ // / _ \/ /
##      \__/\__/_/\_\/_/_/_/\_,_/_//_/_/
##
####### (c) 2020 Institut National de l'Audiovisuel ######################################## Archivage Numérique #######
*/

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"sort"
	"strings"
	"time"

	"github.com/arnumina/eckmuhl.core/pkg/command"
)

const (
	_pluginFuncName = "Export"
)

var (
	_version string
	_builtAt string
)

func findPlugins() (map[string]string, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}

	app := filepath.Base(os.Args[0])

	files, err := filepath.Glob(filepath.Join(filepath.Dir(exe), app+".*.so"))
	if err != nil {
		return nil, err
	}

	plugins := make(map[string]string)

	for _, file := range files {
		plugins[strings.TrimSuffix(strings.TrimPrefix(filepath.Base(file), app+"."), ".so")] = file
	}

	return plugins, nil
}

func cmdHelp() error {
	plugins, err := findPlugins()
	if err != nil {
		return err
	}

	app := filepath.Base(os.Args[0])

	fmt.Println()
	fmt.Println("The command line client")
	fmt.Println("================================================================================")
	fmt.Println("Usage:")
	fmt.Printf("  %s [command [options]]\n", app)
	fmt.Println()
	fmt.Println("Available commands:")

	commands := make([]string, len(plugins)+1)

	n := 0

	for cmd := range plugins {
		commands[n] = cmd
		n++
	}

	commands[n] = "version"

	sort.Strings(commands)

	for _, cmd := range commands {
		fmt.Println("  " + cmd)
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("  Use '%s [command] --help' for more information about a command.\n", app)
	fmt.Println("================================================================================")
	fmt.Println()

	return nil
}

func cmdVersion() {
	builtAt := command.UnixToTime(_builtAt)

	fmt.Println()
	fmt.Println("  eckmuhl")
	fmt.Println("-----------------------------------------------")
	fmt.Println("  version  :", _version)
	fmt.Println("  built at :", builtAt.String())
	fmt.Println("  by       : Archivage Numérique © INA", time.Now().Year())
	fmt.Println("-----------------------------------------------")
	fmt.Println()
}

func runCommand(file string) error {
	plugin, err := plugin.Open(file)
	if err != nil {
		return err
	}

	ef, err := plugin.Lookup(_pluginFuncName)
	if err != nil {
		return err
	}

	fn, ok := ef.(func() command.Command)
	if !ok {
		return fmt.Errorf( /////////////////////////////////////////////////////////////////////////////////////////////
			"this plugin doesn't export the right function: plugin=%s",
			file,
		)
	}

	cmd := fn()

	return cmd.Run(os.Args[2:])
}

func run() error {
	if len(os.Args) == 1 {
		return cmdHelp()
	}

	switch os.Args[1] {
	case "--help", "-help", "help":
		return cmdHelp()
	case "--version", "-version", "version":
		cmdVersion()
		return nil
	}

	plugins, err := findPlugins()
	if err != nil {
		return err
	}

	for name, file := range plugins {
		if name == os.Args[1] {
			return runCommand(file)
		}
	}

	return errors.New("this command does not exist") ///////////////////////////////////////////////////////////////////
}

func main() {
	if err := run(); err != nil {
		if errors.Is(err, command.ErrStopApp) {
			return
		}

		fmt.Fprintf( ///////////////////////////////////////////////////////////////////////////////////////////////////
			os.Stderr,
			"Error: cmd=%s >>> %s\n",
			os.Args[1],
			err,
		)

		os.Exit(1)
	}
}

/*
######################################################################################################## @(°_°)@ #######
*/
