package main

type IScheduler interface {
	Add(p IProcess)
	Next() (IProcess, bool)
	PrintAll()
}

type BasicScheduler struct {
	processes Queue[IProcess]
}

func MakeBasicScheduler() *BasicScheduler {
	s := BasicScheduler{
		processes: Queue[IProcess]{},
	}
	return &s
}

func (s *BasicScheduler) Add(p IProcess) {
	s.processes.Enqueue(p)
}

func (s *BasicScheduler) Next() (IProcess, bool) {
	r, isSuccess := s.processes.Dequeue()
	return r, isSuccess
}

func (s *BasicScheduler) PrintAll() {
	s.processes.PrintAll()
}
