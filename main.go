package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/state-survey/devices"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	workers := 10
	listDevices := devices.CreateDevices(100)
	inputCommands := generateCommands()
	commands := make(chan Command, len(inputCommands))
	errChan := make(chan error)

	var wg sync.WaitGroup
	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, listDevices, commands, errChan)
		}()
	}

	for _, c := range inputCommands {
		commands <- c
	}
	close(commands)

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		fmt.Println("Error: ", err)
		cancel()
	}
}

// simulateDeviceCommand
func simulateDeviceCommand(devices devices.Devices, command Command) error {
	switch command.task {
	case 0:
		err := devices[command.device].Stop()
		if err != nil {
			return err
		}
	case 1:
		err := devices[command.device].Start()
		if err != nil {
			return err
		}
	}

	time.Sleep(time.Millisecond * 10)

	return nil
}

// worker
func worker(ctx context.Context, devices devices.Devices, command <-chan Command, errChan chan<-error) {
	for {
		select {
		case <-ctx.Done():
			errChan <- context.Canceled
			return
		case c, ok := <-command:
			if !ok {
				return
			}
			if err := simulateDeviceCommand(devices, c); err != nil {
				errChan <- err
			}
		}
	}
}


type Command struct {
	device int
	task   int
}

func generateCommands() []Command {
	commands := make([]Command, 0)

	for i := 0; i < 100; i++ {
		commands = append(commands, Command{i, 1})
		commands = append(commands, Command{i, 0})
	}

	return commands
}
