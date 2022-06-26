package main

import (
	"github.com/pzs-pzs/cherry-pick/pkg/flow"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	var repo string
	var out string
	root := &cobra.Command{
		Use:   "cherry-pick analyze",
		Short: "A tool to analyze cherry-pick commit, output yaml file",
	}
	analyze := &cobra.Command{
		Use:   "analyze",
		Short: "A tool to analyze cherry-pick commit, output yaml file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return flow.Run(repo, out)
		},
	}
	root.AddCommand(analyze)

	analyze.PersistentFlags().StringVarP(&repo, "repo", "r", "", "github repo url only https")
	analyze.PersistentFlags().StringVarP(&out, "out", "o", "", "output file path")
	err := root.Execute()
	if err != nil {
		log.Fatalln(err)
		return
	}
}
