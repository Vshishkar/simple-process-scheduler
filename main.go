package main

import "fmt"

func main() {
	CPU := MakeCPU()
	exit := make(chan struct{})
	OS := MakeOs(CPU, exit)

	for i := range 10 {
		OS.CreateProcess(CreateProcessCmd{
			name: fmt.Sprintf("Process %v", i),
		})
	}

	for range exit {
		fmt.Println("exiting...")
	}
}
