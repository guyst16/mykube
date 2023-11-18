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
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/digitalocean/go-libvirt"
	"github.com/dustin/go-humanize"
	"github.com/guyst16/mykube/pkg/embedfiles"
	"github.com/guyst16/mykube/pkg/virtualmachine"
	"github.com/kdomanski/iso9660"
	"github.com/urfave/cli/v2"
)

// Directories
var ASSETS_MYKUBE_DIR = ".mykube"
var ASSETS_MYKUBE_VM_DIR = ""
var LIBVIRT_MYKUBE_DIR = "/var/lib/libvirt/images/mykube"
var LIBVIRT_MYKUBE_UTIL_DIR = LIBVIRT_MYKUBE_DIR + "/" + "util"
var DIRECTORIES_UTIL = [3]string{}

// Utils
// TODO: REPLACE BASE IMAGE URL AND EXTRACT TAR FILE
// var LIBVIRT_MYKUBE_UTIL_BASE_IMAGE_URL = "https://docs.google.com/uc?export=download&confirm=t&id=1yPZ8nk9PeRyek2tRGRkeZVWbaapLbeou"
var LIBVIRT_MYKUBE_UTIL_BASE_IMAGE_URL = "https://github.com/guyst16/mykube/raw/Feat/Refactor_Golang/image-assets/Base-image.qcow2?download="
var LIBVIRT_MYKUBE_UTIL_BASE_IMAGE_PATH = LIBVIRT_MYKUBE_UTIL_DIR + "/" + "Base-image.qcow2"
var LIBVIRT_MYKUBE_VM_DIR = ""
var LIBVIRT_MYKUBE_VM_BASE_IMAGE_PATH = ""
var LIBVIRT_MYKUBE_VM_CLOUDCONFIG_PATH = ""
var LIBVIRT_MYKUBE_VM_METADATA_PATH = ""
var LIBVIRT_MYKUBE_UTIL_CLOUDCONFIG_ISO_PATH = LIBVIRT_MYKUBE_UTIL_DIR + "/" + "cidata.iso"

// Valid
var OS_IMAGE_SHA256SUM = "cafd46df34c9dacb981391e339e00ae582bdcd5d42441bd2708ab54cc5ee856e"

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

					// Get home directory
					userHomeDir, err := os.UserHomeDir()
					if err != nil {
						log.Fatal(err)
					}

					ASSETS_MYKUBE_DIR = userHomeDir + "/" + ASSETS_MYKUBE_DIR
					DIRECTORIES_UTIL = [...]string{ASSETS_MYKUBE_DIR, LIBVIRT_MYKUBE_DIR, LIBVIRT_MYKUBE_UTIL_DIR}

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

					// Create Libvirt vm directory, assets vm directory and copy image
					ASSETS_MYKUBE_VM_DIR = ASSETS_MYKUBE_DIR + "/" + vmName
					LIBVIRT_MYKUBE_VM_DIR = LIBVIRT_MYKUBE_DIR + "/" + vmName
					LIBVIRT_MYKUBE_VM_CLOUDCONFIG_PATH = LIBVIRT_MYKUBE_VM_DIR + "/" + "user-data"
					LIBVIRT_MYKUBE_VM_METADATA_PATH = LIBVIRT_MYKUBE_VM_DIR + "/" + "meta-data"
					LIBVIRT_MYKUBE_VM_BASE_IMAGE_PATH = LIBVIRT_MYKUBE_VM_DIR + "/" + "Base-image.qcow2"
					err = os.Mkdir(LIBVIRT_MYKUBE_VM_DIR, 0744)
					if err != nil {
						log.Fatal(err)
					}
					err = CopyFile(LIBVIRT_MYKUBE_UTIL_BASE_IMAGE_PATH, LIBVIRT_MYKUBE_VM_BASE_IMAGE_PATH)
					if err != nil {
						log.Fatal(err)
					}
					err = os.Mkdir(ASSETS_MYKUBE_VM_DIR, 0744)
					if err != nil {
						log.Fatal(err)
					}

					// Generate ssh key pair inside vm directory
					err = virtualmachine.CreateVirtualmachineSSHKeyPair(ASSETS_MYKUBE_VM_DIR)
					if err != nil {
						log.Fatal(err)
					}

					// Create meta-data files and iso
					cloudConfigContent, _ := embedfiles.InnerReadFile("assets/user-data")
					metaDataContent, _ := embedfiles.InnerReadFile("assets/meta-data")

					vmPubKey, err := os.ReadFile(ASSETS_MYKUBE_VM_DIR + "/public_key.pem")
					if err != nil {
						log.Fatal(err)
					}

					cloudConfigContent = virtualmachine.InjectSSHKeyIntoUserDataYamlFile(cloudConfigContent, string(vmPubKey))

					fmt.Print(string(cloudConfigContent))

					err = os.WriteFile(LIBVIRT_MYKUBE_VM_CLOUDCONFIG_PATH, cloudConfigContent, 0644)
					if err != nil {
						return err
					}

					err = os.WriteFile(LIBVIRT_MYKUBE_VM_METADATA_PATH, metaDataContent, 0644)
					if err != nil {
						return err
					}

					// Create cloud config ISO in VM directory
					cloudConfigFilesArr := []string{LIBVIRT_MYKUBE_VM_CLOUDCONFIG_PATH, LIBVIRT_MYKUBE_VM_METADATA_PATH}
					LIBVIRT_MYKUBE_VM_CLOUDCONFIG_ISO_PATH := LIBVIRT_MYKUBE_VM_DIR + "/" + strings.Split(LIBVIRT_MYKUBE_UTIL_CLOUDCONFIG_ISO_PATH, "/")[len(strings.Split(LIBVIRT_MYKUBE_UTIL_CLOUDCONFIG_ISO_PATH, "/"))-1]
					CreateISO(cloudConfigFilesArr, LIBVIRT_MYKUBE_VM_CLOUDCONFIG_ISO_PATH)

					// Create mykube virtual machine
					myVM := virtualmachine.NewVirtualmachine("os", LIBVIRT_MYKUBE_VM_BASE_IMAGE_PATH, LIBVIRT_MYKUBE_VM_CLOUDCONFIG_ISO_PATH, 1, 1, vmName)
					myVM.CreateVirtualmachine()

					// Start mykube virtual machine
					virtualmachine.StartVirtualMachine(vmName)
					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "delete a single node K8S",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "vm_name",
						Usage:       "Name of the K8S cluster",
						Required:    true,
						Destination: &vmName,
					},
				},
				Action: func(ctx *cli.Context) error {
					err := virtualmachine.DeleteVirtualMachine(vmName)
					if err != nil {
						log.Print(err.Error())
					}

					LIBVIRT_MYKUBE_VM_DIR = LIBVIRT_MYKUBE_DIR + "/" + vmName
					// Get home directory
					userHomeDir, err := os.UserHomeDir()
					if err != nil {
						log.Fatal(err)
					}
					ASSETS_MYKUBE_DIR = userHomeDir + "/" + ASSETS_MYKUBE_DIR
					ASSETS_MYKUBE_VM_DIR = ASSETS_MYKUBE_DIR + "/" + vmName
					os.RemoveAll(LIBVIRT_MYKUBE_VM_DIR)
					os.RemoveAll(ASSETS_MYKUBE_VM_DIR)

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
			{
				Name:  "test",
				Usage: "list all single nodes K8S",
				Action: func(ctx *cli.Context) error {
					addr, err := virtualmachine.GetVirtualMachineIP("test")
					if err != nil {
						log.Fatal(err)
					}
					fmt.Print(addr)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// Progress bar type and functions
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

// Download files to specific destination
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

	// Get file size
	i, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	fileSize := int64(i)

	// create bar
	bar := pb.New(int(fileSize)).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	bar.Start()

	// create proxy reader
	reader := bar.NewProxyReader(resp.Body)

	// and copy from reader
	io.Copy(out, reader)
	bar.Finish()

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

// Create new cidata ISO file
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
