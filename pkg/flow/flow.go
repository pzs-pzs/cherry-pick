package flow

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pkg/errors"
	"github.com/pzs-pzs/cherry-pick/pkg/printer"
	"github.com/pzs-pzs/cherry-pick/pkg/util"
	"github.com/rs/zerolog/log"
	"sync"
)

// NewEngine 初始化处理流程的核心控制器
func NewEngine() *Engine {
	return &Engine{
		worker: make(chan Worker, 100),
		p:      printer.NewYamlPrinter(),
		mutex:  sync.Mutex{},
	}
}

type (
	Engine struct {
		worker chan Worker
		rst    []*CherryCommit
		mutex  sync.Mutex
		p      printer.Printer
		wg     sync.WaitGroup
	}

	CherryCommit struct {
		ID   string
		From string
	}

	PrintData struct {
		OriginalCommit    string   `yaml:"original_commit"`
		CherryPickCommits []string `yaml:"cherry_pick_commits"`
	}
)

// Init 初始化源数据，从远端获取，按照分支的粒度分发到对对应的 Worker
func (e *Engine) Init(url string) error {
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

// Start 并发处理每个分支的cherry-pick日志分析，找出所有的commit封装到 CherryCommit
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

// Output 将所有的 CherryCommit 进行处理，并输出到指定yaml文件
func (e *Engine) Output(path string) error {
	e.wg.Wait()

	cache := make(map[string][]string)
	for _, commit := range e.rst {
		t := cache[commit.From]
		t = append(t, commit.ID)
		cache[commit.From] = t
	}
	var out []*PrintData
	for k, v := range cache {
		out = append(out, &PrintData{
			OriginalCommit:    k,
			CherryPickCommits: util.RemoveDuplication(v),
		})
	}

	err := e.p.Print(path, func() interface{} {
		return out
	})
	if err != nil {
		return errors.WithStack(err)
	}
	log.Info().Msgf("cherry-pick analyze ok! you can get result in [%s]", path)
	return nil
}

func removeDuplication(in []string) (out []string) {
	t := make(map[string]struct{})
	for _, v := range in {
		if _, ok := t[v]; ok {
			continue
		}
		out = append(out, v)
		t[v] = struct{}{}
	}
	return
}

// Run 执行分析cherry-pick的所有组件的flow
func Run(repo, out string) error {
	log.Info().Msgf("repo is [%s] out is [%s]", repo, out)
	err := check(repo, out)
	if err != nil {
		return err
	}
	return run(repo, out)
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
	engine := NewEngine()
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
