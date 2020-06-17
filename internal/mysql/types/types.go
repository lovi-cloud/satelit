package types

import (
	"database/sql/driver"
	"fmt"
	"net"
)

// IPNet is net.IPNet with the implementation of the Valuer and Scanner interface.
type IPNet net.IPNet

// Value implements the database/sql/driver Valuer interface.
func (i IPNet) Value() (driver.Value, error) {
	return driver.Value(i.String()), nil
}

// Scan implements the database/sql Scanner interface.
func (i *IPNet) Scan(src interface{}) error {
	var ipNet *IPNet
	var err error
	switch src := src.(type) {
	case string:
		ipNet, err = parseCIDR(src)
	case []uint8:
		ipNet, err = parseCIDR(fmt.Sprintf("%s", src))
	default:
		return fmt.Errorf("incompatible type for IPNet: %T", src)
	}
	if err != nil {
		return err
	}
	*i = *ipNet
	return nil
}

func (i *IPNet) String() string {
	ipNet := net.IPNet(*i)
	return ipNet.String()
}

// IP is net.IP with the implementation of the Valuer and Scanner interface.
type IP net.IP

// Value implements the database/sql/driver Valuer interface.
func (i IP) Value() (driver.Value, error) {
	return driver.Value(i.String()), nil
}

// Scan implements the database/sql Scanner interface.
func (i *IP) Scan(src interface{}) error {
	var ip *IP
	var err error
	switch src := src.(type) {
	case nil:
		ip = nil
	case string:
		ip, err = parseIP(src)
	case []uint8:
		ip, err = parseIP(fmt.Sprintf("%s", src))
	default:
		return fmt.Errorf("incompatible type for IP: %T", src)
	}
	if err != nil {
		return err
	}
	*i = *ip
	return nil
}

func (i *IP) String() string {
	ip := net.IP(*i)
	return ip.String()
}

func parseCIDR(s string) (*IPNet, error) {
	_, n, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}
	ipNet := IPNet(*n)
	return &ipNet, nil
}

func parseIP(s string) (*IP, error) {
	i := net.ParseIP(s)
	if i == nil {
		return nil, fmt.Errorf("failed to parse IP: input=\"%s\"", i)
	}
	ip := IP(i)
	return &ip, nil
}
