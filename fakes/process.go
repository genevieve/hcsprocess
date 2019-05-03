package fakes

import "io"

type Process struct {
	StdioCall struct {
		Returns struct {
			Stdin  io.WriteCloser
			Stdout io.ReadCloser
			Stderr io.ReadCloser
			Error  error
		}
	}
}

func (p *Process) Stdio() (io.WriteCloser, io.ReadCloser, io.ReadCloser, error) {
	return p.StdioCall.Returns.Stdin, p.StdioCall.Returns.Stdout, p.StdioCall.Returns.Stderr, p.StdioCall.Returns.Error
}

func (p *Process) Wait() error {
	return nil
}

func (p *Process) ExitCode() (int, error) {
	return 0, nil
}
