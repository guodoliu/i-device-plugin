package main

import (
	"github.com/guodoliu/i-device-plugin/pkg/device_plugin"
	"github.com/guodoliu/i-device-plugin/pkg/utils"
	"k8s.io/klog/v2"
)

func main() {
	klog.Infof("device plugin starting")
	dp := device_plugin.NewGopherDevicePlugin()
	go dp.Run()

	// register when device plugin start
	if err := dp.Register(); err != nil {
		klog.Fatalf("register to kubelet failed: %v", err)
	}
	klog.Infof("register to kubelet successfully")

	// watch kubelet.sock, when kubelet restart, exit device plugin, then will restart by DaemonSet
	stop := make(chan struct{})
	err := utils.WatchKubelet(stop)
	if err != nil {
		klog.Fatalf("watch kubelet failed: %v", err)
	}

	<-stop
	klog.Infof("kubelet restart, exiting")
}
