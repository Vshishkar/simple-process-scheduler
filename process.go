package main

import (
	"fmt"
	"time"
)

type CreateProcessCmd struct {
	name string
}

type IProcess interface {
	GetId() int
	NextCommand() (*ProcessCommand, bool)
	MoveToFinished()
	MoveToReady()
	MoveToBlocked()
	MoveToRunning()
	UpdateRunFor(d time.Duration)
	TrySetFirstRunDate(t time.Time) bool
	ToReportItem() ProcessRunReportItem
}

type ProcessStatus string

const (
	running  ProcessStatus = "running"
	ready    ProcessStatus = "ready"
	blocked  ProcessStatus = "blocked"
	finished ProcessStatus = "finished"
)

type ProcessCommandType string

const (
	compute ProcessCommandType = "compute"
	io      ProcessCommandType = "i/o"
)

type ProcessCommand struct {
	args        struct{}
	commandType ProcessCommandType
}

func (c *ProcessCommand) String() string {
	return fmt.Sprintf("cmd: %v", c.commandType)
}

type Process struct {
	id          int
	name        string
	status      ProcessStatus
	runFor      time.Duration
	arrivedAt   time.Time
	finishedAt  time.Time
	commands    []*ProcessCommand
	firstRun    time.Time
	nextCommand int
}

func MakeBasicProcess(id int, name string, arrivedAt time.Time, commands []*ProcessCommand) IProcess {
	p := &Process{
		id:          id,
		name:        name,
		arrivedAt:   arrivedAt,
		commands:    commands,
		nextCommand: 0,
		runFor:      time.Duration(0),
	}
	return p
}

func (p *Process) GetId() int {
	return p.id
}

func (p *Process) NextCommand() (*ProcessCommand, bool) {
	if p.nextCommand >= len(p.commands) {
		return nil, false
	}
	cmd := p.commands[p.nextCommand]
	p.nextCommand += 1

	return cmd, true
}

func (p *Process) MoveToBlocked() {
	p.status = blocked
}

func (p *Process) MoveToFinished() {
	p.status = finished
	p.finishedAt = time.Now()
}

func (p *Process) MoveToReady() {
	p.status = ready
}

func (p *Process) MoveToRunning() {
	p.status = running
}

func (p *Process) String() string {
	return fmt.Sprintf("pid: %v, name %v, aat %v, rnf %v fdt %v", p.id, p.name, p.arrivedAt, p.runFor, p.finishedAt)
}

func (p *Process) UpdateRunFor(d time.Duration) {
	p.runFor += d
}

func (p *Process) TrySetFirstRunDate(t time.Time) bool {
	if !p.firstRun.IsZero() {
		return false
	}

	p.firstRun = t
	return true
}

func (p *Process) ToReportItem() ProcessRunReportItem {
	return ProcessRunReportItem{
		pid:            p.id,
		arrivedAt:      p.arrivedAt,
		finishedAt:     p.finishedAt,
		turnAroundTime: p.finishedAt.Sub(p.arrivedAt),
		tResponseTime:  p.firstRun.Sub(p.arrivedAt),
		runFor:         p.runFor,
	}
}
