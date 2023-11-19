package virtualmachine

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/digitalocean/go-libvirt"
	"github.com/guyst16/mykube/pkg/libvirtconn"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Virtualmachine struct {
	os_name          string
	os_path          string
	cloudconfig_path string
	vcpu_amount      int
	memory_amount    int
	name             string
}

var STATES = map[libvirt.ConnectListAllDomainsFlags]string{libvirt.ConnectListDomainsRunning: "Running", libvirt.ConnectListDomainsPaused: "Paused", libvirt.ConnectListDomainsShutoff: "Shutoff"}

// Create virtual machine object
func NewVirtualmachine(os_name string, os_path string, cloudconfig_path string, vcpu_amount int, memory_amount int, name string) *Virtualmachine {
	v := Virtualmachine{os_name: os_name, os_path: os_path, cloudconfig_path: cloudconfig_path, vcpu_amount: vcpu_amount, memory_amount: memory_amount, name: name}
	return &v
}

func GetVirtualMachine(vmName string) (dom *libvirt.Domain) {
	libvirtconn := libvirtconn.ConnectLibvirtLocal()
	vm, _ := libvirtconn.DomainLookupByName(vmName)
	return &vm
}

func ListAllVirtualmachines() {
	println("ID\tNAME\t\tUUID\t\t\t\t\tSTATE")
	println("-----------------------------------------------------------------------")
	libvirtconn := libvirtconn.ConnectLibvirtLocal()
	for state := range STATES {
		flags := state
		domains, _, _ := libvirtconn.ConnectListAllDomains(1, flags)
		for _, vm := range domains {
			fmt.Printf("%d\t%s\t%x\t%s\n", vm.ID, vm.Name, vm.UUID, STATES[state])
		}
	}
}

func (vm Virtualmachine) CreateVirtualmachine() {
	vmXML := ModifyXML("assets/vmTemplate.xml", vm.name, vm.os_path, vm.cloudconfig_path)
	libvirtconn := libvirtconn.ConnectLibvirtLocal()
	vmXMLString := string(vmXML)
	_, err := libvirtconn.DomainDefineXML(vmXMLString)
	if err != nil {
		log.Fatal(err)
	}
}

// Start defined vm
func StartVirtualMachine(vmName string) {
	libvirtconn := libvirtconn.ConnectLibvirtLocal()
	domain := GetVirtualMachine(vmName)
	if domain == nil {
		log.Fatal("Virtual machine not defined")
	}
	err := libvirtconn.DomainCreate(*domain)
	if err != nil {
		log.Fatal(err)
	}
}

// Delete defined vm
func DeleteVirtualMachine(vmName string) (err error) {
	libvirtconn := libvirtconn.ConnectLibvirtLocal()
	domain := GetVirtualMachine(vmName)
	if domain == nil {
		return errors.New("vm undefined")
	}
	err = libvirtconn.DomainUndefine(*domain)
	if err != nil {
		return err
	}
	err = libvirtconn.DomainDestroy(*domain)
	if err != nil {
		return err
	}

	return nil
}

// Get virtual machine public IP
func GetVirtualMachineIP(vmName string) (vmIPAddress string, err error) {
	libvirtconn := libvirtconn.ConnectLibvirtLocal()
	domain := GetVirtualMachine(vmName)
	domainIPAddress := ""
	if domain == nil {
		return domainIPAddress, errors.New("virtual machine not defined")
	}

	// Get list of the domain interfaces addresses
	intrefacesList, err := libvirtconn.DomainInterfaceAddresses(*domain, 0, 0)
	if err != nil {
		return domainIPAddress, err
	}

	if len(intrefacesList) == 0 {
		return domainIPAddress, errors.New("domain " + vmName + " doesn't have an IP address")
	}

	return intrefacesList[0].Addrs[0].Addr, nil
}

// Copy file from vm to host
func CopyFileFromVirtualMachineToHost(sshClient *ssh.Client, remoteFilePath string, localFilePath string) (err error) {
	defer sshClient.Close()

	// Open an SFTP session on the remote machine
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		fmt.Println("Failed to create SFTP client:", err)
		return
	}
	defer sftpClient.Close()

	// Check if the file exists
	_, err = sftpClient.Stat("/etc/kubernetes/admin.conf")
	if err != nil {
		return err
	}

	// Copy file to local
	// File exists, open the remote file
	remoteFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		fmt.Println("Error opening remote file:", err)
		return
	}
	defer remoteFile.Close()

	// Create a local file to copy the contents to
	localFile, err := os.Create(localFilePath)
	if err != nil {
		fmt.Println("Error creating local file:", err)
		return
	}
	defer localFile.Close()

	// Copy the contents of the remote file to the local file
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return
	}

	return nil
}

// Get ssh client for a virtaul machine
func GetVirtualMachineSSHConnection(vmName string, vmPubKeyPath string) (shhClient *ssh.Client, err error) {
	privateKey, err := LoadPrivateKeyFromFile(vmPubKeyPath)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: "sumit",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privateKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	vmIP, err := GetVirtualMachineIP(vmName)
	if err != nil {
		return nil, err
	}

	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", vmIP, 22), config)
	if err != nil {
		return nil, err
	}

	return sshClient, nil
}

func LoadPrivateKeyFromFile(filePath string) (ssh.Signer, error) {
	// Read the contents of the file
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(fileContent)
	if err != nil {
		return nil, err
	}

	return key, nil
}
