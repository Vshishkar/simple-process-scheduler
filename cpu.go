package main

import (
	"fmt"
	"time"
)

type CPU struct {
	processes chan *ProcessTrap
	cycle     time.Duration
}

type ProcessTrap struct {
	process       IProcess
	interruptTrap chan ProcessExecutionResult
}

type ProcessExecutionResult struct {
	startTime  time.Time
	pid        int
	runFor     time.Duration
	isFinished bool
}

func (p *ProcessExecutionResult) GetId() int {
	return p.pid
}

func MakeCPU() *CPU {
	c := &CPU{
		processes: make(chan *ProcessTrap),
		cycle:     time.Second * 5,
	}

	c.run()
	return c
}

func (c *CPU) run() {
	go func() {
		for {
			trap := <-c.processes
			fmt.Printf("CPU: running %v process \n", trap.process.GetId())
			startTime := time.Now()

			timer := time.NewTimer(c.cycle)
			defer timer.Stop()

			isFinished := false
		loop:
			for {
				timer.Reset(c.cycle)
				select {
				case <-timer.C:
					break loop
				default:
					{
						cmd, success := trap.process.NextCommand()
						if !success {
							isFinished = true
							break loop
						}
						fmt.Printf("CPU: pid: %v executing command %v \n", trap.process.GetId(), cmd)
						time.Sleep(time.Millisecond * 125)
					}
				}
			}

			finished := time.Now()
			fmt.Printf("CPU: finished running %v process \n", trap.process.GetId())
			trap.interruptTrap <- ProcessExecutionResult{
				pid:        trap.process.GetId(),
				runFor:     finished.Sub(startTime),
				isFinished: isFinished,
				startTime:  startTime,
			}
		}
	}()
}
