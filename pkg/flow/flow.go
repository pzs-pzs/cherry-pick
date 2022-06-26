package flow

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pkg/errors"
	"github.com/pzs-pzs/cherry-pick/pkg/printer"
	"github.com/rs/zerolog/log"
	"sync"
)

func NewEngine() *Engine {
	return &Engine{
		worker: make(chan Worker, 100),
		p:      printer.NewYamlPrinter(),
		mutex:  sync.Mutex{},
	}
}

type Engine struct {
	worker chan Worker
	rst    []*CherryCommit
	mutex  sync.Mutex
	p      printer.Printer
	wg     sync.WaitGroup
}

type CherryCommit struct {
	ID   string
	From string
}

func (e *Engine) Init(url string) error {
	if url == "" {
		return errors.New("invalid url,url is empty")
	}
	log.Info().Msgf("start clone [%s] to local", url)
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	log.Info().Msgf("clone [%s] to local is ok", url)
	references, err := r.References()
	if err != nil {
		return errors.WithStack(err)
	}
	var remoteRefs []*plumbing.Reference
	err = references.ForEach(func(ref *plumbing.Reference) error {
		if !ref.Name().IsRemote() {
			return nil
		}
		remoteRefs = append(remoteRefs, ref)
		return nil
	})
	if err != nil {
		return errors.WithStack(err)
	}
	for _, ref := range remoteRefs {
		e.worker <- &worker{
			ref:  ref,
			repo: r,
		}
		e.wg.Add(1)
	}
	close(e.worker)
	return nil
}

func (e *Engine) Start() {

	for {
		select {
		case w := <-e.worker:
			if w == nil {
				log.Info().Msgf("worker channel closed")
				return
			}
			go func() {
				defer e.wg.Done()
				out, err := w.Do()
				if err != nil {
					return
				}
				e.mutex.Lock()
				defer e.mutex.Unlock()
				e.rst = append(e.rst, out...)
			}()
		}
	}

}

func (e *Engine) OutPut(path string) error {
	e.wg.Wait()

	cache := make(map[string][]string)
	for _, commit := range e.rst {
		log.Info().Msgf("%+v", commit)
		t := cache[commit.From]
		t = append(t, commit.ID)
		cache[commit.From] = t
	}
	var out []*PrintData
	for k, v := range cache {
		out = append(out, &PrintData{
			OriginalCommit:    k,
			CherryPickCommits: v,
		})
	}

	return e.p.Print(path, func() interface{} {
		return out
	})
}

type PrintData struct {
	OriginalCommit    string   `yaml:"original_commit"`
	CherryPickCommits []string `yaml:"cherry_pick_commits"`
}
