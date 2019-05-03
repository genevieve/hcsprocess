package hcsprocess

import (
	"io"
	"sync"
	"time"
)

type Process struct {
	process process
}

type process interface {
	Stdio() (io.WriteCloser, io.ReadCloser, io.ReadCloser, error)
	Wait() error
	ExitCode() (int, error)
	CloseStdin()
}

func New(p process) *Process {
	return &Process{process: p}
}

func (p *Process) AttachIO(attachStdin io.Reader, attachStdout, attachStderr io.Writer) (int, error) {
	stdin, stdout, stderr, err := p.process.Stdio()
	if err != nil {
		return -1, err
	}

	var wg sync.WaitGroup

	if attachStdin != nil {
		wg.Add(1)
		go func() {
			_, _ = io.Copy(stdin, attachStdin)
			_ = stdin.Close()
			p.process.CloseStdin()
			wg.Done()
		}()
	} else {
		_ = stdin.Close()
		// p.process.CloseStdin()
	}

	if attachStdout != nil {
		wg.Add(1)
		go func() {
			_, _ = io.Copy(attachStdout, stdout)
			_ = stdout.Close()
			wg.Done()
		}()
	} else {
		_ = stdout.Close()
	}

	if attachStderr != nil {
		wg.Add(1)
		go func() {
			_, _ = io.Copy(attachStderr, stderr)
			_ = stderr.Close()
			wg.Done()
		}()
	} else {
		_ = stderr.Close()
	}

	err = p.process.Wait()
	waitWithTimeout(&wg, 1*time.Second)
	if err != nil {
		return -1, err
	}

	return p.process.ExitCode()
}

func waitWithTimeout(wg *sync.WaitGroup, timeout time.Duration) {
	wgEmpty := make(chan interface{}, 1)
	go func() {
		wg.Wait()
		wgEmpty <- nil
	}()

	select {
	case <-time.After(timeout):
	case <-wgEmpty:
	}
}
