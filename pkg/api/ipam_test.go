package api

import (
	"testing"

	"github.com/go-test/deep"

	pb "github.com/whywaita/satelit/api/satelit"
)

func TestSatelitServer_CreateSubnet(t *testing.T) {
	_, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	tests := []struct {
		input *pb.CreateSubnetRequest
		want  *pb.CreateSubnetResponse
		err   bool
	}{
		{
			input: &pb.CreateSubnetRequest{
				Name:           "subnet1000",
				Network:        "192.0.2.0/24",
				VlanId:         1000,
				Start:          "192.0.2.100",
				End:            "192.0.2.200",
				Gateway:        "192.0.2.1",
				DnsServer:      "8.8.8.8",
				MetadataServer: "192.0.2.15",
			},
			want: &pb.CreateSubnetResponse{
				Subnet: &pb.Subnet{
					Uuid:           "",
					Name:           "subnet1000",
					Network:        "192.0.2.0/24",
					VlanId:         1000,
					Start:          "192.0.2.100",
					End:            "192.0.2.200",
					Gateway:        "192.0.2.1",
					DnsServer:      "8.8.8.8",
					MetadataServer: "192.0.2.15",
				},
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.CreateSubnet(ctx, test.input)
		if got != nil {
			test.want.Subnet.Uuid = got.Subnet.Uuid
		}
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_GetSubnet(t *testing.T) {
	_, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	subnetResp, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "subnet1000",
		Network:        "192.0.2.0/24",
		VlanId:         1000,
		Start:          "192.0.2.100",
		End:            "192.0.2.200",
		Gateway:        "192.0.2.1",
		DnsServer:      "8.8.8.8",
		MetadataServer: "192.0.2.15",
	})
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}

	tests := []struct {
		input *pb.GetSubnetRequest
		want  *pb.GetSubnetResponse
		err   bool
	}{
		{
			input: &pb.GetSubnetRequest{
				Uuid: subnetResp.Subnet.Uuid,
			},
			want: &pb.GetSubnetResponse{
				Subnet: subnetResp.Subnet,
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.GetSubnet(ctx, test.input)
		if got != nil {
			test.want.Subnet.Uuid = got.Subnet.Uuid
		}
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_ListSubnet(t *testing.T) {
	_, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	subnetResp, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "subnet1000",
		Network:        "192.0.2.0/24",
		VlanId:         1000,
		Start:          "192.0.2.100",
		End:            "192.0.2.200",
		Gateway:        "192.0.2.1",
		DnsServer:      "8.8.8.8",
		MetadataServer: "192.0.2.15",
	})
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}

	tests := []struct {
		input *pb.ListSubnetRequest
		want  *pb.ListSubnetResponse
		err   bool
	}{
		{
			input: &pb.ListSubnetRequest{},
			want: &pb.ListSubnetResponse{
				Subnets: []*pb.Subnet{
					subnetResp.Subnet,
				},
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.ListSubnet(ctx, test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_DeleteSubnet(t *testing.T) {
	_, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	subnetResp, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "subnet1000",
		Network:        "192.0.2.0/24",
		VlanId:         1000,
		Start:          "192.0.2.100",
		End:            "192.0.2.200",
		Gateway:        "192.0.2.1",
		DnsServer:      "8.8.8.8",
		MetadataServer: "192.0.2.15",
	})
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}

	tests := []struct {
		input *pb.DeleteSubnetRequest
		want  *pb.DeleteSubnetResponse
		err   bool
	}{
		{
			input: &pb.DeleteSubnetRequest{
				Uuid: subnetResp.Subnet.Uuid,
			},
			want: &pb.DeleteSubnetResponse{},
			err:  false,
		},
	}
	for _, test := range tests {
		got, err := client.DeleteSubnet(ctx, test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_CreateAddress(t *testing.T) {
	_, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	subnetResp, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "subnet1000",
		Network:        "192.0.2.0/24",
		VlanId:         1000,
		Start:          "192.0.2.100",
		End:            "192.0.2.200",
		Gateway:        "192.0.2.1",
		DnsServer:      "8.8.8.8",
		MetadataServer: "192.0.2.15",
	})
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}

	tests := []struct {
		input *pb.CreateAddressRequest
		want  *pb.CreateAddressResponse
		err   bool
	}{
		{
			input: &pb.CreateAddressRequest{
				SubnetId: subnetResp.Subnet.Uuid,
			},
			want: &pb.CreateAddressResponse{
				Address: &pb.Address{
					Uuid:     "",
					Ip:       "192.0.2.100",
					SubnetId: subnetResp.Subnet.Uuid,
				},
			},
			err: false,
		},
		{
			input: &pb.CreateAddressRequest{
				SubnetId: subnetResp.Subnet.Uuid,
				FixedIp:  "192.0.2.110",
			},
			want: &pb.CreateAddressResponse{
				Address: &pb.Address{
					Uuid:     "",
					Ip:       "192.0.2.110",
					SubnetId: subnetResp.Subnet.Uuid,
				},
			},
			err: false,
		},
		{
			input: &pb.CreateAddressRequest{
				SubnetId: subnetResp.Subnet.Uuid,
				FixedIp:  "192.0.2.120",
			},
			want: &pb.CreateAddressResponse{
				Address: &pb.Address{
					Uuid:     "",
					Ip:       "192.0.2.120",
					SubnetId: subnetResp.Subnet.Uuid,
				},
			},
			err: false,
		},
		{
			input: &pb.CreateAddressRequest{
				SubnetId: subnetResp.Subnet.Uuid,
			},
			want: &pb.CreateAddressResponse{
				Address: &pb.Address{
					Uuid:     "",
					Ip:       "192.0.2.101",
					SubnetId: subnetResp.Subnet.Uuid,
				},
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.CreateAddress(ctx, test.input)
		if got != nil {
			test.want.Address.Uuid = got.Address.Uuid
		}
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_GetAddress(t *testing.T) {
	_, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	subnetResp, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "subnet1000",
		Network:        "192.0.2.0/24",
		VlanId:         1000,
		Start:          "192.0.2.100",
		End:            "192.0.2.200",
		Gateway:        "192.0.2.1",
		DnsServer:      "8.8.8.8",
		MetadataServer: "192.0.2.15",
	})
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}
	addressResp, err := client.CreateAddress(ctx, &pb.CreateAddressRequest{
		SubnetId: subnetResp.Subnet.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test address: %+v", err)
	}

	tests := []struct {
		input *pb.GetAddressRequest
		want  *pb.GetAddressResponse
		err   bool
	}{
		{
			input: &pb.GetAddressRequest{
				Uuid: addressResp.Address.Uuid,
			},
			want: &pb.GetAddressResponse{
				Address: addressResp.Address,
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.GetAddress(ctx, test.input)
		if got != nil {
			test.want.Address.Uuid = got.Address.Uuid
		}
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_ListAddress(t *testing.T) {
	_, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	subnetResp, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "subnet1000",
		Network:        "192.0.2.0/24",
		VlanId:         1000,
		Start:          "192.0.2.100",
		End:            "192.0.2.200",
		Gateway:        "192.0.2.1",
		DnsServer:      "8.8.8.8",
		MetadataServer: "192.0.2.15",
	})
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}
	addressResp, err := client.CreateAddress(ctx, &pb.CreateAddressRequest{
		SubnetId: subnetResp.Subnet.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test address: %+v", err)
	}

	tests := []struct {
		input *pb.ListAddressRequest
		want  *pb.ListAddressResponse
		err   bool
	}{
		{
			input: &pb.ListAddressRequest{
				SubnetId: subnetResp.Subnet.Uuid,
			},
			want: &pb.ListAddressResponse{
				Addresses: []*pb.Address{
					addressResp.Address,
				},
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.ListAddress(ctx, test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_DeleteAddress(t *testing.T) {
	_, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	subnetResp, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "subnet1000",
		Network:        "192.0.2.0/24",
		VlanId:         1000,
		Start:          "192.0.2.100",
		End:            "192.0.2.200",
		Gateway:        "192.0.2.1",
		DnsServer:      "8.8.8.8",
		MetadataServer: "192.0.2.15",
	})
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}
	addressResp, err := client.CreateAddress(ctx, &pb.CreateAddressRequest{
		SubnetId: subnetResp.Subnet.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test address: %+v", err)
	}

	tests := []struct {
		input *pb.DeleteAddressRequest
		want  *pb.DeleteAddressResponse
		err   bool
	}{
		{
			input: &pb.DeleteAddressRequest{
				Uuid: addressResp.Address.Uuid,
			},
			want: &pb.DeleteAddressResponse{},
			err:  false,
		},
	}
	for _, test := range tests {
		got, err := client.DeleteAddress(ctx, test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_CreateLease(t *testing.T) {
	_, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	subnetResp, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "subnet1000",
		Network:        "192.0.2.0/24",
		VlanId:         1000,
		Start:          "192.0.2.100",
		End:            "192.0.2.200",
		Gateway:        "192.0.2.1",
		DnsServer:      "8.8.8.8",
		MetadataServer: "192.0.2.15",
	})
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}
	addressResp, err := client.CreateAddress(ctx, &pb.CreateAddressRequest{
		SubnetId: subnetResp.Subnet.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test address: %+v", err)
	}

	tests := []struct {
		input *pb.CreateLeaseRequest
		want  *pb.CreateLeaseResponse
		err   bool
	}{
		{
			input: &pb.CreateLeaseRequest{
				AddressId: addressResp.Address.Uuid,
			},
			want: &pb.CreateLeaseResponse{
				Lease: &pb.Lease{
					Uuid:       "",
					MacAddress: "",
					AddressId:  addressResp.Address.Uuid,
				},
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.CreateLease(ctx, test.input)
		if got != nil {
			test.want.Lease.Uuid = got.Lease.Uuid
			test.want.Lease.MacAddress = got.Lease.MacAddress
		}
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_GetLease(t *testing.T) {
	_, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	subnetResp, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "subnet1000",
		Network:        "192.0.2.0/24",
		VlanId:         1000,
		Start:          "192.0.2.100",
		End:            "192.0.2.200",
		Gateway:        "192.0.2.1",
		DnsServer:      "8.8.8.8",
		MetadataServer: "192.0.2.15",
	})
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}
	addressResp, err := client.CreateAddress(ctx, &pb.CreateAddressRequest{
		SubnetId: subnetResp.Subnet.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test address: %+v", err)
	}
	leaseResp, err := client.CreateLease(ctx, &pb.CreateLeaseRequest{
		AddressId: addressResp.Address.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test lease: %+v", err)
	}

	tests := []struct {
		input *pb.GetLeaseRequest
		want  *pb.GetLeaseResponse
		err   bool
	}{
		{
			input: &pb.GetLeaseRequest{
				Uuid: leaseResp.Lease.Uuid,
			},
			want: &pb.GetLeaseResponse{
				Lease: leaseResp.Lease,
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.GetLease(ctx, test.input)
		if got != nil {
			test.want.Lease.Uuid = got.Lease.Uuid
			test.want.Lease.MacAddress = got.Lease.MacAddress
		}
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_ListLease(t *testing.T) {
	_, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	subnetResp, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "subnet1000",
		Network:        "192.0.2.0/24",
		VlanId:         1000,
		Start:          "192.0.2.100",
		End:            "192.0.2.200",
		Gateway:        "192.0.2.1",
		DnsServer:      "8.8.8.8",
		MetadataServer: "192.0.2.15",
	})
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}
	addressResp, err := client.CreateAddress(ctx, &pb.CreateAddressRequest{
		SubnetId: subnetResp.Subnet.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test address: %+v", err)
	}
	leaseResp, err := client.CreateLease(ctx, &pb.CreateLeaseRequest{
		AddressId: addressResp.Address.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test lease: %+v", err)
	}

	tests := []struct {
		input *pb.ListLeaseRequest
		want  *pb.ListLeaseResponse
		err   bool
	}{
		{
			input: &pb.ListLeaseRequest{},
			want: &pb.ListLeaseResponse{
				Leases: []*pb.Lease{
					leaseResp.Lease,
				},
			},
			err: false,
		},
	}
	for _, test := range tests {
		got, err := client.ListLease(ctx, test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}

func TestSatelitServer_DeleteLease(t *testing.T) {
	_, teardownTeleskop, err := setupTeleskop(nil)
	if err != nil {
		t.Fatalf("failed to get teleskop endpoint %+v\n", err)
	}
	defer teardownTeleskop()

	ctx, client, teardown := getSatelitClient()
	defer teardown()

	subnetResp, err := client.CreateSubnet(ctx, &pb.CreateSubnetRequest{
		Name:           "subnet1000",
		Network:        "192.0.2.0/24",
		VlanId:         1000,
		Start:          "192.0.2.100",
		End:            "192.0.2.200",
		Gateway:        "192.0.2.1",
		DnsServer:      "8.8.8.8",
		MetadataServer: "192.0.2.15",
	})
	if err != nil {
		t.Fatalf("failed to create test subnet: %+v", err)
	}
	addressResp, err := client.CreateAddress(ctx, &pb.CreateAddressRequest{
		SubnetId: subnetResp.Subnet.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test address: %+v", err)
	}
	leaseResp, err := client.CreateLease(ctx, &pb.CreateLeaseRequest{
		AddressId: addressResp.Address.Uuid,
	})
	if err != nil {
		t.Fatalf("failed to create test lease: %+v", err)
	}

	tests := []struct {
		input *pb.DeleteLeaseRequest
		want  *pb.DeleteLeaseResponse
		err   bool
	}{
		{
			input: &pb.DeleteLeaseRequest{
				Uuid: leaseResp.Lease.Uuid,
			},
			want: &pb.DeleteLeaseResponse{},
			err:  false,
		},
	}
	for _, test := range tests {
		got, err := client.DeleteLease(ctx, test.input)
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}
		if diff := deep.Equal(test.want, got); len(diff) != 0 {
			t.Fatalf("want %q, but %q, diff %q:", test.want, got, diff)
		}
	}
}
