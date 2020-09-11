package api

import (
	"bytes"
	"context"
	"fmt"

	"github.com/whywaita/satelit/pkg/europa"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/whywaita/satelit/api/satelit"
	"github.com/whywaita/satelit/internal/logger"
	"github.com/whywaita/satelit/internal/qcow2"
)

// ShowVolume call GetVolume to Europa Backend
func (s *SatelitServer) ShowVolume(ctx context.Context, req *pb.ShowVolumeRequest) (*pb.ShowVolumeResponse, error) {
	v, err := s.Datastore.GetVolume(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve volume from datastore: %+v", err)
	}
	e, ok := s.Europa[v.BackendName]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "invalid backend name (name: %s)", v.BackendName)
	}

	volume, err := e.GetVolume(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve volume: %+v", err)
	}

	return &pb.ShowVolumeResponse{
		Volume: volume.ToPb(),
	}, nil
}

// ListVolume call ListVolume to Europa Backend
func (s *SatelitServer) ListVolume(ctx context.Context, req *pb.ListVolumeRequest) (*pb.ListVolumeResponse, error) {
	var volumes []europa.Volume
	for _, e := range s.Europa {
		vs, err := e.ListVolume(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to retrieve volumes: %+v", err)
		}

		volumes = append(volumes, vs...)
	}

	var pvs []*pb.Volume
	for _, v := range volumes {
		pvs = append(pvs, v.ToPb())
	}

	return &pb.ListVolumeResponse{
		Volumes: pvs,
	}, nil
}

// AddVolume call CreateVolume to Europa backend
func (s *SatelitServer) AddVolume(ctx context.Context, req *pb.AddVolumeRequest) (*pb.AddVolumeResponse, error) {
	u, err := s.parseRequestUUID(req.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request name (need uuid): %+v", err)
	}
	e, ok := s.Europa[req.BackendName]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid backend name (name: %s)", req.BackendName)
	}

	volume, err := e.CreateVolume(ctx, u, int(req.CapacityGigabyte))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create volume: %+v", err)
	}

	return &pb.AddVolumeResponse{
		Volume: volume.ToPb(),
	}, nil
}

// AddVolumeImage call CreateVolumeImage to Europa backend
func (s *SatelitServer) AddVolumeImage(ctx context.Context, req *pb.AddVolumeImageRequest) (*pb.AddVolumeImageResponse, error) {
	u, err := s.parseRequestUUID(req.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request name (need uuid): %+v", err)
	}
	e, ok := s.Europa[req.BackendName]
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid backend name (name: %s)", req.BackendName)
	}

	sourceImageID, err := s.parseRequestUUID(req.SourceImageId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request source image id (need uuid): %+v", err)
	}

	v, err := e.CreateVolumeFromImage(ctx, u, int(req.CapacityGigabyte), sourceImageID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create volume from image: %+v", err)
	}

	return &pb.AddVolumeImageResponse{
		Volume: v.ToPb(),
	}, nil
}

// AttachVolume call AttachVolume to Europa backend
func (s *SatelitServer) AttachVolume(ctx context.Context, req *pb.AttachVolumeRequest) (*pb.AttachVolumeResponse, error) {
	v, err := s.Datastore.GetVolume(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve volume from datastore: %+v", err)
	}
	e, ok := s.Europa[v.BackendName]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "invalid backend name (name: %s)", v.BackendName)
	}

	if _, _, err := e.AttachVolumeTeleskop(ctx, req.Id, req.Hostname); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to attach volume: %+v", err)
	}

	return &pb.AttachVolumeResponse{}, nil
}

// DetachVolume call DetachVolume to Europa backend
func (s *SatelitServer) DetachVolume(ctx context.Context, req *pb.DetachVolumeRequest) (*pb.DetachVolumeResponse, error) {
	v, err := s.Datastore.GetVolume(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve volume from datastore: %+v", err)
	}
	e, ok := s.Europa[v.BackendName]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "invalid backend name (name: %s)", v.BackendName)
	}

	if err := e.DetachVolume(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to detach volume: %+v", err)
	}

	return &pb.DetachVolumeResponse{}, nil
}

// DeleteVolume call DeleteVolume to Europa backend
func (s *SatelitServer) DeleteVolume(ctx context.Context, req *pb.DeleteVolumeRequest) (*pb.DeleteVolumeResponse, error) {
	v, err := s.Datastore.GetVolume(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve volume from datastore: %+v", err)
	}
	e, ok := s.Europa[v.BackendName]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "invalid backend name (name: %s)", v.BackendName)
	}

	if err := e.DeleteVolume(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete volume: %+v", err)
	}

	return &pb.DeleteVolumeResponse{}, nil
}

// ListImage retrieves all images
func (s *SatelitServer) ListImage(ctx context.Context, req *pb.ListImageRequest) (*pb.ListImageResponse, error) {
	var images []europa.BaseImage

	for _, e := range s.Europa {
		is, err := e.ListImage()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to retrieves images: %+v", err)
		}

		images = append(images, is...)
	}

	var pbImages []*pb.Image
	for _, image := range images {
		pbImages = append(pbImages, image.ToPb())
	}

	return &pb.ListImageResponse{
		Images: pbImages,
	}, nil
}

// UploadImage upload to europa backend
func (s *SatelitServer) UploadImage(stream pb.Satelit_UploadImageServer) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	buf := pool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		pool.Put(buf)
	}()

	m, err := s.receiveImage(stream, buf)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to receive image file: %+v", err)
	}
	logger.Logger.Debug(fmt.Sprintf("received image (name: %s)", m.name))

	e, ok := s.Europa[m.europaBackendName]
	if !ok {
		return status.Errorf(codes.InvalidArgument, "invalid backend name (name: %s)", m.europaBackendName)
	}

	// validate qcow2 image
	b := buf.Bytes()
	reader := bytes.NewReader(b)
	isQcow2, header := qcow2.Probe(reader)
	if isQcow2 == false {
		return status.Errorf(codes.InvalidArgument, "failed to validate qcow2 image")
	}

	// send to europa
	image, err := e.UploadImage(ctx, b, m.name, m.description, sanitizeImageSize(header.Size))
	if err != nil {
		return status.Errorf(codes.Internal, "failed to upload image to europa: %+v", err)
	}
	logger.Logger.Debug("uploaded image to europa")
	defer func() {
		uuid := image.UUID
		if err != nil {
			if err := e.DeleteImage(ctx, uuid); err != nil {
				logger.Logger.Warn(fmt.Sprint("failed to delete image: %w", err))
			}
		}
	}()

	err = stream.SendAndClose(&pb.UploadImageResponse{Image: image.ToPb()})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to send and close: %+v", err)
	}
	logger.Logger.Debug("close UploadImage stream")

	// save to image info in database
	err = s.Datastore.PutImage(*image)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to put image: %+v", err)
	}
	logger.Logger.Debug("completed write image to datastore")

	logger.Logger.Info(fmt.Sprintf("UploadImage is successfully! (name: %s)", m.name))
	return nil
}

// DeleteImage delete image
func (s *SatelitServer) DeleteImage(ctx context.Context, req *pb.DeleteImageRequest) (*pb.DeleteImageResponse, error) {
	imageID, err := s.parseRequestUUID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to parse request image id (need uuid): %+v", err)
	}
	image, err := s.Datastore.GetImage(imageID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve image from datastore: %+v", err)
	}
	v, err := s.Datastore.GetVolume(ctx, image.CacheVolumeID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve cache volume from datastore: %+v", err)
	}
	e, ok := s.Europa[v.BackendName]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "invalid backend name (name: %s)", v.BackendName)
	}

	err = e.DeleteImage(ctx, imageID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete image from europa: %+v", err)
	}

	return &pb.DeleteImageResponse{}, nil
}
