package virtualmachine

import (
	"fmt"

	"github.com/digitalocean/go-libvirt"
	"github.com/guyst16/mykube/mykubeLibvirt"
)

type Virtualmachine struct {
	os_name       string
	os_path       string
	vcpu_amount   int
	memory_amount int
	name          string
}

// Create virtual machine object
func NewVirtualmachine(os_name string, os_path string, vcpu_amount int, memory_amount int, name string) *Virtualmachine {
	v := Virtualmachine{os_name: os_name, os_path: os_path, vcpu_amount: vcpu_amount, memory_amount: memory_amount, name: name}
	return &v
}

func ListAllVirtualmachines() {
	println("ID\tNAME\t\tUUID\t\t\t\t\tSTATE")
	println("-----------------------------------------------------------------------")
	libvirtconn := mykubeLibvirt.ConnectLibvirtLocal()
	states := map[libvirt.ConnectListAllDomainsFlags]string{libvirt.ConnectListDomainsRunning: "Running", libvirt.ConnectListDomainsPaused: "Paused", libvirt.ConnectListDomainsShutoff: "Shutoff"}
	for state := range states {
		flags := state
		domains, _, _ := libvirtconn.ConnectListAllDomains(1, flags)
		for _, vm := range domains {
			fmt.Printf("%d\t%s\t%x\t%s\n", vm.ID, vm.Name, vm.UUID, states[state])
			print(libvirtconn.DomainGetXMLDesc(vm, 1))
		}
	}
}

func (vm Virtualmachine) CreateVirtualmachine() {
	// libvirtconn := mykubeLibvirt.ConnectLibvirtLocal()
}
