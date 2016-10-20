package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type IApplication interface {
	Run()
	IsDebug() bool
	GetResource() string
	GetCwd() string
}

type CommandLineHandler func(c *CommandLine) error

type CommandLine struct {
	Root *cobra.Command
	Main CommandLineHandler

	debug    bool
	resource string
	cwd      string
}

func (c *CommandLine) GetResource() string {
	return c.resource
}

func (c *CommandLine) GetCwd() string {
	return c.cwd
}

func (c *CommandLine) IsDebug() bool {
	return c.debug
}

func (c *CommandLine) Run() {

	var err error

	c.resource = filepath.Dir(os.Args[0])
	c.cwd, err = os.Getwd()

	if err != nil {
		println(err.Error())
		return
	}

	err = c.Main(c)

	if(err != nil){
		println(err.Error())
		return
	}

	c.Root.PersistentPreRun = func(_ *cobra.Command, _ []string) {
		if c.IsDebug() {
			fmt.Printf("Run in debug mode...\n")
			fmt.Printf("debug:\t%v\n", c.debug)
			fmt.Printf("resource:\t%v\n", c.resource)
			fmt.Printf("cwd:\t%v\n", c.cwd)
		}
	}

	flags := c.Root.PersistentFlags()

	flags.BoolVar(&c.debug, "debug", false, "run in debug model")
	flags.StringVar(&c.resource, "data-dir", c.resource, "app resource directory")

	c.Root.Execute()

}

func (c *CommandLine) AddCommand(cmds ...*cobra.Command) {
	c.Root.AddCommand(cmds...)
}
