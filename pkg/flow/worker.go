package flow

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
	"strings"
)

type Worker interface {
	Do() ([]*CherryCommit, error)
}

type worker struct {
	ref     *plumbing.Reference
	repo    *git.Repository
	commits []*CherryCommit
}

func (w *worker) Do() ([]*CherryCommit, error) {
	l, err := w.repo.Log(&git.LogOptions{
		From: w.ref.Hash(),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = l.ForEach(func(c *object.Commit) error {
		return w.analyze(c)
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return w.commits, nil
}

// analyze 分析cherry-pick commit
// example msg '(cherry picked from commit 8fe8f231cf539e3346a4fd31d9c275bf168f6cc8)'
func (w *worker) analyze(c *object.Commit) error {
	if !strings.Contains(c.Message, "cherry picked from commit") {
		return nil
	}
	index := strings.Index(c.Message, "cherry picked from commit")

	from := c.Message[index : len(c.Message)-2]
	split := strings.Split(from, " ")
	if len(split) < 1 {
		return nil
	}
	w.commits = append(w.commits, &CherryCommit{
		ID:   c.Hash.String(),
		From: split[len(split)-1],
	})
	return nil
}

func (w *worker) Out() []*CherryCommit {
	return w.commits
}
