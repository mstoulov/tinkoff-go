package storage

import (
	"context"
	"sync/atomic"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

func (r *Result) AtomicAdd(other Result) {
	atomic.AddInt64(&r.Count, other.Count)
	atomic.AddInt64(&r.Size, other.Size)
}

func (r *Result) Add(other Result) {
	r.Count += other.Count
	r.Size += other.Size
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	//maxWorkersCount int
	//sem             semaphore.Weighted
	cancel context.CancelCauseFunc
	ctx    context.Context
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	//return &sizer{maxWorkersCount: 1000, sem: *semaphore.NewWeighted(int64(1000))}
	return &sizer{}
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	a.ctx, a.cancel = context.WithCancelCause(ctx)
	resCh := make(chan Result)
	go a.RecursiveSize(a.ctx, d, resCh)
	failed := false
	for {
		select {
		case <-a.ctx.Done():
			failed = true
		case result := <-resCh:
			if failed {
				return Result{}, context.Cause(a.ctx)
			} else {
				return result, nil
			}
		}
	}
}

func (a *sizer) RecursiveSize(parCtx context.Context, d Dir, parResCh chan Result) {
	if parCtx.Err() != nil {
		parResCh <- Result{}
		return
	}

	ctx, cancel := context.WithCancel(parCtx)
	_ = cancel
	dirs, files, lsErr := d.Ls(ctx)
	//_ = lsErr
	if lsErr != nil {
		a.cancel(lsErr)
		parResCh <- Result{}
		return
	}

	if parCtx.Err() != nil {
		parResCh <- Result{}
		return
	}

	result := Result{Count: int64(len(files)), Size: 0}
	resCh := make(chan Result)
	for _, dir := range dirs {
		//a.sem.Acquire(ctx, 1)
		go a.RecursiveSize(ctx, dir, resCh)
	}

	for _, file := range files {
		size, statErr := file.Stat(ctx)
		//_ = statErr
		if statErr != nil {
			a.cancel(statErr)
			break
		}
		result.Size += size
	}

	interrupted := false
	for i := 0; i < len(dirs); i++ {
		select {
		case <-ctx.Done():
			i--
			interrupted = true
			parResCh <- Result{}
		case res := <-resCh:
			result.Add(res)
		}
	}
	if !interrupted {
		parResCh <- result
	}
	//a.sem.Release(1)
}
