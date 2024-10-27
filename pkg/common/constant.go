package common

import "time"

const (
	ResourceName   string = "gordon.com/gopher"
	DevicePath     string = "/etc/gophers"
	DeviceSocket   string = "gopher.sock"
	ConnectTimeout        = time.Second * 5
)
