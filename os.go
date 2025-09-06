package main

import (
	"fmt"
	"time"
)

type OS struct {
	scheduler        IScheduler
	createProcessCmd chan CreateProcessCmd
	cpu              *CPU
	pidCounter       int
	processMap       map[int]IProcess
	exit             chan struct{}
}

func (o *OS) CreateProcess(cmd CreateProcessCmd) {
	o.createProcessCmd <- cmd
}

func MakeOs(cpu *CPU, exit chan struct{}) *OS {
	os := OS{}
	os.createProcessCmd = make(chan CreateProcessCmd)
	os.processMap = make(map[int]IProcess)
	os.cpu = cpu
	os.scheduler = MakeBasicScheduler()
	os.pidCounter = 0
	os.exit = exit

	os.run()
	os.createProcessListener()

	return &os
}

func (o *OS) createProcessListener() {
	go func() {
		for cmd := range o.createProcessCmd {
			fmt.Printf("OS: handing create process command. name: %v \n", cmd.name)
			p := MakeBasicProcess(o.pidCounter, cmd.name, time.Now(), GenerateCommands(10))
			o.processMap[p.GetId()] = p
			o.pidCounter += 1
			o.scheduler.Add(p)
		}
	}()
}

func GenerateCommands(length int) []*ProcessCommand {
	cmds := make([]*ProcessCommand, length)
	for i := range length {
		cmds[i] = &ProcessCommand{
			commandType: compute,
		}
	}
	return cmds
}

type ProcessRunReportItem struct {
	pid            int
	arrivedAt      time.Time
	finishedAt     time.Time
	runFor         time.Duration
	turnAroundTime time.Duration
	tResponseTime  time.Duration
}

type ProcessRunReport struct {
	items                 []ProcessRunReportItem
	averageTurnAroundTime time.Duration
	averageResponseTime   time.Duration
}

func (p ProcessRunReport) String() string {
	return fmt.Sprintf("pCount %v, avg turnaround %v, avg resp %v", len(p.items), p.averageTurnAroundTime, p.averageResponseTime)
}

func (o *OS) createProcessRunReport() ProcessRunReport {
	pCount := len(o.processMap)
	items := make([]ProcessRunReportItem, pCount)

	counter := 0
	for _, p := range o.processMap {
		items[counter] = p.ToReportItem()
	}

	turnAroundTimes := make([]time.Duration, pCount)
	responseTimes := make([]time.Duration, pCount)
	for i, report := range items {
		turnAroundTimes[i] = report.turnAroundTime
		responseTimes[i] = report.tResponseTime
	}

	return ProcessRunReport{
		items:                 items,
		averageTurnAroundTime: averageDuration(turnAroundTimes),
		averageResponseTime:   averageDuration(responseTimes),
	}
}

func (o *OS) run() {
	go func() {
		for {
			processToRun, isSuccess := o.scheduler.Next()

			if !isSuccess {
				report := o.createProcessRunReport()
				fmt.Println("report: ")
				fmt.Println(report)
				close(o.exit)
				return
			}

			trap := ProcessTrap{
				process:       processToRun,
				interruptTrap: make(chan ProcessExecutionResult),
			}

			fmt.Printf("OS: picked process to run. %v \n", processToRun)
			o.cpu.processes <- &trap
			executionResult := <-trap.interruptTrap

			processToRun.UpdateRunFor(executionResult.runFor)
			if executionResult.isFinished {
				processToRun.MoveToFinished()
			} else {
				processToRun.MoveToReady()
				o.scheduler.Add(processToRun)
			}

			processToRun.TrySetFirstRunDate(executionResult.startTime)
			fmt.Printf("OS: finished process execution. p: %v \n", processToRun)
		}
	}()
}
