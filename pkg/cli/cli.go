package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/google/uuid"
	"github.com/guyst16/mykube/pkg/embedfiles"
	"github.com/guyst16/mykube/pkg/virtualmachine"
	"github.com/kdomanski/iso9660"
	"github.com/urfave/cli/v2"
)

// Directories
var MAIN_DIR = ".mykube"
var LIBVIRT_MYKUBE_DIR = "/var/lib/libvirt/images/mykube"
var LIBVIRT_MYKUBE_UTIL_DIR = LIBVIRT_MYKUBE_DIR + "/" + "util"
var DIRECTORIES_UTIL = [3]string{}

// Utils
var LIBVIRT_MYKUBE_UTIL_BASE_IMAGE_URL = "https://download.fedoraproject.org/pub/fedora/linux/releases/38/Cloud/x86_64/images/Fedora-Cloud-Base-38-1.6.x86_64.qcow2"
var LIBVIRT_MYKUBE_UTIL_BASE_IMAGE_PATH = LIBVIRT_MYKUBE_UTIL_DIR + "/" + "Base-image.qcow2"
var LIBVIRT_MYKUBE_VM_DIR = ""
var LIBVIRT_MYKUBE_VM_BASE_IMAGE_PATH = ""
var LIBVIRT_MYKUBE_VM_CLOUDCONFIG_ISO_PATH = ""
var LIBVIRT_MYKUBE_UTIL_CLOUDCONFIG_PATH = LIBVIRT_MYKUBE_UTIL_DIR + "/" + "user-data"
var LIBVIRT_MYKUBE_UTIL_METADATA_PATH = LIBVIRT_MYKUBE_UTIL_DIR + "/" + "meta-data"
var LIBVIRT_MYKUBE_UTIL_CLOUDCONFIG_ISO_PATH = LIBVIRT_MYKUBE_UTIL_DIR + "/" + "cidata.iso"

// Valid
var OS_IMAGE_SHA256SUM = "d334670401ff3d5b4129fcc662cf64f5a6e568228af59076cc449a4945318482"

// TODO: Handle all err inside the helpers functions
func Cli() {
	var vmName string

	app := &cli.App{
		Name:     "mykube",
		Usage:    "Manage single node K8S",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Guy Steinberger",
				Email: "guyst16@gmail.com",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "Create a single node K8S",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "vm_name",
						Usage:       "Name of the K8S cluster",
						Required:    true,
						Destination: &vmName,
					},
				},
				Action: func(ctx *cli.Context) error {
					// Validate vm doesn't already exist
					_, err := GetVM(vmName)
					if err == nil {
						log.Fatalf("Virtual machine named %q already exists", vmName)
					}

					// Create home directory
					userHomeDir, err := os.UserHomeDir()
					if err != nil {
						log.Fatal(err)
					}

					MAIN_DIR = userHomeDir + "/" + MAIN_DIR
					DIRECTORIES_UTIL = [...]string{MAIN_DIR, LIBVIRT_MYKUBE_DIR, LIBVIRT_MYKUBE_UTIL_DIR}

					// Validate directories existence
					for _, dir := range DIRECTORIES_UTIL {
						_, dir_err := os.Stat(dir)
						if os.IsNotExist(dir_err) {
							err := os.Mkdir(dir, 0744)
							if err != nil {
								log.Fatal(err)
							}
						}
					}

					err = ValidateOSImage(LIBVIRT_MYKUBE_UTIL_BASE_IMAGE_PATH, OS_IMAGE_SHA256SUM)
					if err != nil {
						// Download cloud base image
						log.Print("Download OS image file")
						DownloadFile(LIBVIRT_MYKUBE_UTIL_BASE_IMAGE_PATH, LIBVIRT_MYKUBE_UTIL_BASE_IMAGE_URL)
					}

					vm_uuid := uuid.New()

					// Create VM directory and copy image
					LIBVIRT_MYKUBE_VM_DIR = LIBVIRT_MYKUBE_DIR + "/" + vm_uuid.String()
					LIBVIRT_MYKUBE_VM_BASE_IMAGE_PATH = LIBVIRT_MYKUBE_VM_DIR + "/" + "Base-image.qcow2"
					err = os.Mkdir(LIBVIRT_MYKUBE_VM_DIR, 0744)
					if err != nil {
						log.Fatal(err)
					}
					err = CopyFile(LIBVIRT_MYKUBE_UTIL_BASE_IMAGE_PATH, LIBVIRT_MYKUBE_VM_BASE_IMAGE_PATH)
					if err != nil {
						log.Fatal(err)
					}

					// Create meta-data files and iso
					cloudConfigContent, _ := embedfiles.InnerReadFile("assets/user-data")
					metaDataContent, _ := embedfiles.InnerReadFile("assets/meta-data")

					err = os.WriteFile(LIBVIRT_MYKUBE_UTIL_CLOUDCONFIG_PATH, cloudConfigContent, 0644)
					if err != nil {
						return err
					}

					err = os.WriteFile(LIBVIRT_MYKUBE_UTIL_METADATA_PATH, metaDataContent, 0644)
					if err != nil {
						return err
					}

					cloudConfigFilesArr := []string{LIBVIRT_MYKUBE_UTIL_CLOUDCONFIG_PATH, LIBVIRT_MYKUBE_UTIL_METADATA_PATH}
					CreateISO(cloudConfigFilesArr, LIBVIRT_MYKUBE_UTIL_CLOUDCONFIG_ISO_PATH)

					// Copy cloud config ISO to VM directory
					LIBVIRT_MYKUBE_VM_CLOUDCONFIG_ISO_PATH = LIBVIRT_MYKUBE_VM_DIR + "/" + strings.Split(LIBVIRT_MYKUBE_UTIL_CLOUDCONFIG_ISO_PATH, "/")[len(strings.Split(LIBVIRT_MYKUBE_UTIL_CLOUDCONFIG_ISO_PATH, "/"))-1]
					CopyFile(LIBVIRT_MYKUBE_UTIL_CLOUDCONFIG_ISO_PATH, LIBVIRT_MYKUBE_VM_CLOUDCONFIG_ISO_PATH)

					// Create mykube virtual machine
					myVM := virtualmachine.NewVirtualmachine("os", LIBVIRT_MYKUBE_VM_BASE_IMAGE_PATH, LIBVIRT_MYKUBE_VM_CLOUDCONFIG_ISO_PATH, 1, 1, vmName)
					myVM.CreateVirtualmachine()

					// Start mykube virtual machine
					virtualmachine.StartVM(vmName)
					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "delete a single node K8S",
				Action: func(ctx *cli.Context) error {
					fmt.Printf("delete %s\n", ctx.Args().Get(0))
					return nil
				},
			},
			{
				Name:  "connect",
				Usage: "connect a single node K8S",
				Action: func(ctx *cli.Context) error {
					fmt.Printf("connect %s\n", ctx.Args().Get(0))
					return nil
				},
			},
			{
				Name:  "list",
				Usage: "list all single nodes K8S",
				Action: func(ctx *cli.Context) error {
					virtualmachine.ListAllVirtualmachines()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func DownloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func CopyFile(src string, dest string) (err error) {
	bytesRead, err := os.ReadFile(src)

	if err != nil {
		return err
	}

	err = os.WriteFile(dest, bytesRead, 0644)

	if err != nil {
		return err
	}

	return err
}

func CreateISO(filesListFullPath []string, outputISOPath string) {
	writer, err := iso9660.NewWriter()
	if err != nil {
		log.Fatalf("failed to create writer: %s", err)
	}
	defer writer.Cleanup()

	for _, element := range filesListFullPath {
		f, err := os.Open(element)
		if err != nil {
			log.Fatalf("failed to open file: %s", err)
		}
		defer f.Close()

		err = writer.AddFile(f, strings.Split(element, "/")[len(strings.Split(element, "/"))-1])
		if err != nil {
			log.Fatalf("failed to add file: %s", err)
		}
	}

	outputFile, err := os.OpenFile(outputISOPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("failed to create file: %s", err)
	}

	err = writer.WriteTo(outputFile, "cidata")
	if err != nil {
		log.Fatalf("failed to write ISO image: %s", err)
	}

	err = outputFile.Close()
	if err != nil {
		log.Fatalf("failed to close output file: %s", err)
	}
}

// Validate the given image file, else return an error
func ValidateOSImage(filePath string, validSHA256sum string) (err error) {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return err
	}

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	imageSum := h.Sum(nil)

	if hex.EncodeToString(imageSum) != validSHA256sum {
		return errors.New("SHA256sum is not equal")
	}

	return nil
}

// Return vm details
func GetVM(vmName string) (vm *libvirt.Domain, err error) {
	vm = virtualmachine.GetVirtualMachine(vmName)
	if vm.Name == "" {
		return nil, errors.New("virtual machine does not exist")
	}

	return vm, nil
}
