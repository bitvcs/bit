package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

type job struct {
	kind     jobKind
	password string
	hash     string
	result   chan jobResult
}

type jobKind int

const (
	hashJob jobKind = iota
	compareJob
)

type jobResult struct {
	hash    string
	err     error
	matches bool
}

type Hasher struct {
	jobs    chan job
	workers []*worker
}

type worker struct {
	h    *Hasher
	quit chan struct{}
}

func NewHasher(workerCount int) *Hasher {
	if workerCount < 1 {
		workerCount = 1
	}
	h := &Hasher{
		jobs: make(chan job),
	}
	for i := 0; i < workerCount; i++ {
		w := &worker{h: h, quit: make(chan struct{})}
		h.workers = append(h.workers, w)
		go w.run()
	}
	return h
}

func (w *worker) run() {
	for {
		select {
		case j := <-w.h.jobs:
			w.h.execute(j)
		case <-w.quit:
			return
		}
	}
}

func (h *Hasher) execute(j job) {
	switch j.kind {
	case hashJob:
		hash, err := bcrypt.GenerateFromPassword([]byte(j.password), bcrypt.DefaultCost)
		j.result <- jobResult{hash: string(hash), err: err}
	case compareJob:
		err := bcrypt.CompareHashAndPassword([]byte(j.hash), []byte(j.password))
		j.result <- jobResult{matches: err == nil}
	}
}

func (h *Hasher) Hash(password string) (string, error) {
	result := make(chan jobResult, 1)
	h.jobs <- job{kind: hashJob, password: password, result: result}
	res := <-result
	return res.hash, res.err
}

func (h *Hasher) Compare(hash, password string) bool {
	result := make(chan jobResult, 1)
	h.jobs <- job{kind: compareJob, hash: hash, password: password, result: result}
	res := <-result
	return res.matches
}

func (h *Hasher) Close() {
	for _, w := range h.workers {
		close(w.quit)
	}
}
