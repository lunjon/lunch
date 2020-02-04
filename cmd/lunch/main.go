package main

import (
	"fmt"
	"github.com/lunjon/lunch/internal/pkg/edison"
	"github.com/lunjon/lunch/internal/pkg/lunch"
	"github.com/lunjon/lunch/internal/pkg/menu"
	"github.com/lunjon/lunch/internal/pkg/pieplow"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
)

const (
	verboseFlag = "verbose"
	todayFlag   = "today"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "lunch",
		Short: "Get the menu for the local restaurants",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			verbose, _ := cmd.Flags().GetBool(verboseFlag)
			if !verbose {
				log.SetOutput(ioutil.Discard)
			}
		},
	}

	rootCmd.PersistentFlags().BoolP(verboseFlag, "v", false, "Adds debug logs to output")
	rootCmd.PersistentFlags().BoolP(todayFlag, "t", false, "Only output today's courses")

	edisonCmd := &cobra.Command{
		Use:     "edison",
		Aliases: []string{"e", "ed"},
		Short:   "Get the menu for Edison this week",
		Run:     run,
	}

	pieplowCmd := &cobra.Command{
		Use:     "pieplow",
		Aliases: []string{"p", "pi", "pie"},
		Short:   "Get the menu for Pieplow lunch this week",
		Run:     run,
	}

	rootCmd.AddCommand(edisonCmd, pieplowCmd)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, _ []string) {
	var l lunch.Lunch
	switch cmd.Name() {
	case "edison":
		l = edison.New()
	case "pieplow":
		l = pieplow.New()
	default:
		panic("unknown command")
	}

	m, err := l.GetMenu()
	if err != nil {
		fmt.Printf("Failed to get m from %s: %v\n", m.Name(), err)
		os.Exit(1)
	}
	today, err := cmd.Flags().GetBool(todayFlag)

	if today {
		m.FilterDay(menu.Today)
	}

	m.Render()
}
