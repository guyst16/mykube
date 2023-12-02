package libvirtconn

import (
	"log"

	"github.com/digitalocean/go-libvirt"
	"github.com/digitalocean/go-libvirt/socket/dialers"
)

func ConnectLibvirtLocal() *libvirt.Libvirt {
	// Connect to Libvirt unix domain socket
	localSocket := dialers.NewLocal()
	libvirtconn := libvirt.NewWithDialer(localSocket)
	err := libvirtconn.Connect()
	if err != nil {
		log.Fatalf("Failed to dial libvirt socket: %v", err)
	}

	return libvirtconn
}
