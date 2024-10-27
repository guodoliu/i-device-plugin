package device_plugin

import (
	"context"
	"github.com/guodoliu/i-device-plugin/pkg/common"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog/v2"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"net"
	"os"
	"path"
	"syscall"
	"time"
)

type GopherDevicePlugin struct {
	server *grpc.Server
	stop   chan struct{}
	dm     *DeviceMonitor
}

func NewGopherDevicePlugin() *GopherDevicePlugin {
	return &GopherDevicePlugin{
		server: grpc.NewServer(grpc.EmptyServerOption{}),
		stop:   make(chan struct{}),
		dm:     NewDeviceMonitor(common.DevicePath),
	}
}

// Run start gRPC server and watcher
func (g *GopherDevicePlugin) Run() error {
	err := g.dm.List()
	if err != nil {
		klog.Fatalf("failed to list devices: %v", err)
	}

	go func() {
		if err = g.dm.Watch(); err != nil {
			klog.Infof("failed to watch devices: %v", err)
		}
	}()

	pluginapi.RegisterDevicePluginServer(g.server, g)
	// delete old unix socket before start
	socket := path.Join(pluginapi.DevicePluginPath, common.DeviceSocket)
	err = syscall.Unlink(socket)
	if err != nil && !os.IsNotExist(err) {
		return errors.WithMessagef(err, "delete socket %s failed", socket)
	}

	sock, err := net.Listen("unix", socket)
	if err != nil {
		return errors.WithMessagef(err, "listen unix %s failed", socket)
	}

	go g.server.Serve(sock)

	// Wait for server to start by launching a blocking connection
	conn, err := connect(common.DeviceSocket, 5*time.Second)
	if err != nil {
		return err
	}
	conn.Close()

	return nil
}

func connect(socket string, timeout time.Duration) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	c, err := grpc.DialContext(ctx, socket,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			if deadline, ok := ctx.Deadline(); ok {
				return net.DialTimeout("unix", s, time.Until(deadline))
			}
			return net.DialTimeout("unix", s, common.ConnectTimeout)
		}))
	if err != nil {
		return nil, err
	}

	return c, nil
}
