package devices

import "fmt"

type IDevice interface {
	Start() error
	Stop() error
}

func New(name int) IDevice {
	return &Device{name: name}
}

type Device struct {
	name int
}

type Devices map[int]IDevice

func (d *Device) Start() error {
	fmt.Printf("device %v started\n", d.name)
	return nil
}

func (d *Device) Stop() error {
	fmt.Printf("device %v stopped\n", d.name)
	return nil
}

func CreateDevices(count int) Devices {
	devices := make(Devices, count)

	for i := 0; i < count; i++ {
		devices[i] = New(i)
	}

	return devices
}