package main

import (
	"github.com/pzs-pzs/cherry-pick/pkg/flow"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	var repo string
	var out string
	c := &cobra.Command{
		Use:   "pick",
		Short: "A tool to analyze cherry-pick commit, output yaml file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(repo, out)
		},
	}
	c.PersistentFlags().StringVarP(&repo, "repo", "r", "", "github repo url only https")
	c.PersistentFlags().StringVarP(&out, "out", "o", "", "output file path")
	err := c.Execute()
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func run(repo, out string) error {
	engine := flow.NewEngine()
	err := engine.Init(repo)
	if err != nil {
		return err
	}
	go func() {
		engine.Start()
	}()

	err = engine.OutPut(out)

	if err != nil {
		return err
	}
	return nil
}
