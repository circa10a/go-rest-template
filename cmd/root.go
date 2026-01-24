package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	project      = "go-rest-template"
	envVarPrefix = "APP"

	rootCmd = &cobra.Command{
		Use:   project,
		Short: "A template project for Go REST API's",
	}
)

// define config opts to be used by cobra + viber for configuration
type flagDef struct {
	Name      string
	Shorthand string
	Type      string // "bool", "string", "stringArray", "int"
	Default   interface{}
	Usage     string
	ViperKey  string
}

// RegisterFlagTypes registers flags on the provided cobra command according
// to the provided definitions.
func RegisterFlagTypes(cmd *cobra.Command, defs []flagDef) {
	for _, d := range defs {
		switch d.Type {
		case "bool":
			cmd.Flags().BoolP(d.Name, d.Shorthand, d.Default.(bool), d.Usage)
		case "string":
			cmd.Flags().StringP(d.Name, d.Shorthand, d.Default.(string), d.Usage)
		case "stringArray":
			cmd.Flags().StringArrayP(d.Name, d.Shorthand, d.Default.([]string), d.Usage)
		case "int":
			cmd.Flags().IntP(d.Name, d.Shorthand, d.Default.(int), d.Usage)
		}
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
