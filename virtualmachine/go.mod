module github.com/guyst16/mykube/virtualmachine

go 1.19

require (
    "github.com/guyst16/mykube/mykubeLibvirt" v0.0.0
)

replace "github.com/guyst16/mykube/mykubeLibvirt" v0.0.0 => ../libvirt

