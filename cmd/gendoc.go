package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/dhax/go-base/api"
	"github.com/go-chi/docgen"
	"github.com/spf13/cobra"
)

var (
	routes bool
)

// gendocCmd represents the gendoc command
var gendocCmd = &cobra.Command{
	Use:   "gendoc",
	Short: "Generate project documentation",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if routes {
			genRoutesDoc()
		}
	},
}

func init() {
	RootCmd.AddCommand(gendocCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gendocCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gendocCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	gendocCmd.Flags().BoolVarP(&routes, "routes", "r", false, "create api routes markdown file")
}

func genRoutesDoc() {
	api, _ := api.New(false)
	fmt.Print("generating routes markdown file: ")
	md := docgen.MarkdownRoutesDoc(api, docgen.MarkdownOpts{
		ProjectPath: "github.com/dhax/go-base",
		Intro:       "GoBase REST API.",
	})
	if err := ioutil.WriteFile("routes.md", []byte(md), 0644); err != nil {
		log.Println(err)
		return
	}
	fmt.Println("OK")
}
