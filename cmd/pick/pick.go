package main

import (
	"github.com/pkg/errors"
	"github.com/pzs-pzs/cherry-pick/pkg/flow"
	"github.com/pzs-pzs/cherry-pick/pkg/util"
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
		Use:   "cherry-pick analyze",
		Short: "A tool to analyze cherry-pick commit, output yaml file",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := check(repo, out)
			if err != nil {
				return err
			}
			return run(repo, out)
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

func check(repo, out string) error {
	if repo == "" {
		return errors.New("invalid repo,repo is empty")
	}
	if out == "" {
		return errors.New("path is empty, plz check")
	}

	exists, err := util.PathExists(out)
	if err != nil {
		return err
	}
	if exists {
		return errors.Errorf("[%s] already exist", out)
	}
	return nil
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

	err = engine.Output(out)
	if err != nil {
		return err
	}
	return nil
}
