package europa

import pb "github.com/whywaita/satelit/api/satelit"

// A Volume is volume information
type Volume struct {
	ID          string
	Attached    bool
	HostName    string
	CapacityGB  uint32
	BaseImageID string
}

// ToPb parse to Satelit API Server pb.Volume
func (v *Volume) ToPb() *pb.Volume {
	pv := &pb.Volume{
		Id:           v.ID,
		Attached:     v.Attached,
		Hostname:     v.HostName,
		CapacityByte: v.CapacityGB,
	}
	return pv
}

// BaseImage is os Image
type BaseImage struct {
	ID            string
	BaseOS        string
	VolumeID      string // volume id (create by storage backend)
	CacheVolumeID string // if CacheVolumeID is blank, not have cache image
}
