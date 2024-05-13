package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/state-survey/devices"
)

// Мок устройства
type MockDevice struct {
	mock.Mock
}

func (m *MockDevice) Start() error {
	return m.Called().Error(0)
}

func (m *MockDevice) Stop() error {
	return m.Called().Error(0)
}

func TestSimulateDeviceCommand(t *testing.T) {
	mockDevice := new(MockDevice)

	startErr := errors.New("error from start")
	stopErr := errors.New("error from stop")

	mockDevice.On("Start").Return(startErr)
	mockDevice.On("Stop").Return(stopErr)

	devices := make(devices.Devices, 1)
	devices[1] = mockDevice

	testCases := []struct {
		desc	string
		command Command
		want error
	}{
		{
			desc: "start error",
			command: Command{1, 1},
			want: startErr,
		},
		{
			desc: "stop error",
			command: Command{1, 0},
			want: stopErr,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := simulateDeviceCommand(devices, tC.command)
			if err == nil {
				t.Errorf("%s: expected error > %s , got no error", tC.desc,tC.want)
			}
		})
	}

	mockDevice.AssertExpectations(t)
}

func TestContexCancel(t *testing.T)  {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	devices := make(devices.Devices, 1)
	devices[1] = &MockDevice{}
	commands := make(chan Command)
	errChan := make(chan error)

	go worker(ctx, devices, commands, errChan)

	go func() {
		time.Sleep(time.Second * 1)
		cancel()
	}()

	for err := range errChan {
		if !errors.Is(err, context.Canceled) {
			t.Errorf("expect context cancel, got: %v", err)
		} else {
			return
		}
	}
}

