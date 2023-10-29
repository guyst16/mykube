package virtualmachine

import (
	"encoding/xml"
	"log"
	"strings"

	"github.com/guyst16/mykube/pkg/embedfiles"
)

type rawXML struct {
	Inner []byte     `xml:",innerxml"`
	Attrs []xml.Attr `xml:",any,attr"`
}

type diskSubTag struct {
	Attrs []xml.Attr `xml:",any,attr"`
}

type source struct {
	File string `xml:"file,attr"`
}

type disk struct {
	Driver  diskSubTag `xml:"driver"`
	Source  source     `xml:"source"`
	Target  diskSubTag `xml:"target"`
	Address diskSubTag `xml:"address"`
}

type devices struct {
	Emulator  rawXML `xml:"emulator"`
	Disk      []disk `xml:"disk"`
	Interface rawXML `xml:"interface"`
	Serial    rawXML `xml:"serial"`
	Console   rawXML `xml:"console"`
	Channel   rawXML `xml:"channel"`
	Graphics  rawXML `xml:"graphics"`
}

type domain struct {
	Type        string  `xml:"type,attr"`
	Name        string  `xml:"name"`
	Memory      rawXML  `xml:"memory"`
	Vcpu        rawXML  `xml:"vcpu"`
	OS          rawXML  `xml:"os"`
	Features    rawXML  `xml:"features"`
	Cpu         rawXML  `xml:"cpu"`
	Clock       rawXML  `xml:"clock"`
	On_poweroff rawXML  `xml:"on_poweroff"`
	On_reboot   rawXML  `xml:"on_reboot"`
	On_crash    rawXML  `xml:"on_crash"`
	PM          rawXML  `xml:"pm"`
	Devices     devices `xml:"devices"`
}

func ModifyXML(assetFilePath string, vmName string, vmBaseImagePath string, vmCloudConfigIsoPath string) (output []byte) {
	data, err := embedfiles.InnerReadFile(assetFilePath)
	if err != nil {
		log.Fatal(err)
	}

	var domain domain
	if err := xml.Unmarshal(data, &domain); err != nil {
		log.Fatal(err)
	}

	// Modify vm values
	// Modify vm name
	domain.Name = vmName

	// Modify OS disk path and cloud config disk path
	for i := 0; i < len(domain.Devices.Disk); i++ {
		if strings.Contains(domain.Devices.Disk[i].Source.File, "Fedora") {
			domain.Devices.Disk[i].Source.File = vmBaseImagePath
		} else if strings.Contains(domain.Devices.Disk[i].Source.File, "cidata") {
			domain.Devices.Disk[i].Source.File = vmCloudConfigIsoPath
		}
	}

	modified, err := xml.MarshalIndent(&domain, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return modified
}
