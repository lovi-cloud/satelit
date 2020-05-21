package europa

import pb "github.com/whywaita/satelit/api/satelit"

// A Volume is volume information
type Volume struct {
	ID          string
	Attached    bool
	HostName    string
	CapacityGB  uint32
	BaseImageID string
	HostLUNID   int // set when attached
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
	Name          string
	Description   string
	CacheVolumeID string // if CacheVolumeID is blank, not have cache image
}

// ToPb parse to Satelit API Server pb.Image
func (bi *BaseImage) ToPb() *pb.Image {
	i := &pb.Image{
		Id:          bi.ID,
		Name:        bi.Name,
		VolumeId:    bi.CacheVolumeID,
		Description: bi.Description,
	}

	return i
}
