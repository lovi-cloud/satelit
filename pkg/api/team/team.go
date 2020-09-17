package team

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

var mutex = &sync.Mutex{}
var teamMap = map[uint]*net.IPNet{}

// GetTeamID retrieves team id according to the address given
func GetTeamID(address string) (uint, error) {
	mutex.Lock()
	defer mutex.Unlock()

	addr := net.ParseIP(address)
	if addr == nil {
		return 0, fmt.Errorf("failed to parse request address: %s", address)
	}
	for id, ipNet := range teamMap {
		if ipNet.Contains(addr) {
			return id, nil
		}
	}
	return 0, fmt.Errorf("failed to find team id")
}

func init() {
	subnets := strings.Split(subnetList, "\n")
	for i, subnet := range subnets {
		_, ipNet, err := net.ParseCIDR(subnet)
		if err != nil {
			panic(fmt.Errorf("failed to init team library"))
		}
		teamMap[uint(i)+1] = ipNet
	}
}

const subnetList = `10.160.1.0/24
10.160.2.0/24
10.160.3.0/24
10.160.4.0/24
10.160.5.0/24
10.160.6.0/24
10.160.7.0/24
10.160.8.0/24
10.160.9.0/24
10.160.10.0/24
10.160.11.0/24
10.160.12.0/24
10.160.13.0/24
10.160.14.0/24
10.160.15.0/24
10.160.16.0/24
10.160.17.0/24
10.160.18.0/24
10.160.19.0/24
10.160.20.0/24
10.160.21.0/24
10.160.22.0/24
10.160.23.0/24
10.160.24.0/24
10.160.25.0/24
10.160.26.0/24
10.160.27.0/24
10.160.28.0/24
10.160.29.0/24
10.160.30.0/24
10.160.31.0/24
10.160.32.0/24
10.160.33.0/24
10.160.34.0/24
10.160.35.0/24
10.160.36.0/24
10.160.37.0/24
10.160.38.0/24
10.160.39.0/24
10.160.40.0/24
10.160.41.0/24
10.160.42.0/24
10.160.43.0/24
10.160.44.0/24
10.160.45.0/24
10.160.46.0/24
10.160.47.0/24
10.160.48.0/24
10.160.49.0/24
10.160.50.0/24
10.160.51.0/24
10.160.52.0/24
10.160.53.0/24
10.160.54.0/24
10.160.55.0/24
10.160.56.0/24
10.160.57.0/24
10.160.58.0/24
10.160.59.0/24
10.160.60.0/24
10.160.61.0/24
10.160.62.0/24
10.160.63.0/24
10.160.64.0/24
10.160.65.0/24
10.160.66.0/24
10.160.67.0/24
10.160.68.0/24
10.160.69.0/24
10.160.70.0/24
10.160.71.0/24
10.160.72.0/24
10.160.73.0/24
10.160.74.0/24
10.160.75.0/24
10.160.76.0/24
10.160.77.0/24
10.160.78.0/24
10.160.79.0/24
10.160.80.0/24
10.160.81.0/24
10.160.82.0/24
10.160.83.0/24
10.160.84.0/24
10.160.85.0/24
10.160.86.0/24
10.160.87.0/24
10.160.88.0/24
10.160.89.0/24
10.160.90.0/24
10.160.91.0/24
10.160.92.0/24
10.160.93.0/24
10.160.94.0/24
10.160.95.0/24
10.160.96.0/24
10.160.97.0/24
10.160.98.0/24
10.160.99.0/24
10.161.0.0/24
10.161.1.0/24
10.161.2.0/24
10.161.3.0/24
10.161.4.0/24
10.161.5.0/24
10.161.6.0/24
10.161.7.0/24
10.161.8.0/24
10.161.9.0/24
10.161.10.0/24
10.161.11.0/24
10.161.12.0/24
10.161.13.0/24
10.161.14.0/24
10.161.15.0/24
10.161.16.0/24
10.161.17.0/24
10.161.18.0/24
10.161.19.0/24
10.161.20.0/24
10.161.21.0/24
10.161.22.0/24
10.161.23.0/24
10.161.24.0/24
10.161.25.0/24
10.161.26.0/24
10.161.27.0/24
10.161.28.0/24
10.161.29.0/24
10.161.30.0/24
10.161.31.0/24
10.161.32.0/24
10.161.33.0/24
10.161.34.0/24
10.161.35.0/24
10.161.36.0/24
10.161.37.0/24
10.161.38.0/24
10.161.39.0/24
10.161.40.0/24
10.161.41.0/24
10.161.42.0/24
10.161.43.0/24
10.161.44.0/24
10.161.45.0/24
10.161.46.0/24
10.161.47.0/24
10.161.48.0/24
10.161.49.0/24
10.161.50.0/24
10.161.51.0/24
10.161.52.0/24
10.161.53.0/24
10.161.54.0/24
10.161.55.0/24
10.161.56.0/24
10.161.57.0/24
10.161.58.0/24
10.161.59.0/24
10.161.60.0/24
10.161.61.0/24
10.161.62.0/24
10.161.63.0/24
10.161.64.0/24
10.161.65.0/24
10.161.66.0/24
10.161.67.0/24
10.161.68.0/24
10.161.69.0/24
10.161.70.0/24
10.161.71.0/24
10.161.72.0/24
10.161.73.0/24
10.161.74.0/24
10.161.75.0/24
10.161.76.0/24
10.161.77.0/24
10.161.78.0/24
10.161.79.0/24
10.161.80.0/24
10.161.81.0/24
10.161.82.0/24
10.161.83.0/24
10.161.84.0/24
10.161.85.0/24
10.161.86.0/24
10.161.87.0/24
10.161.88.0/24
10.161.89.0/24
10.161.90.0/24
10.161.91.0/24
10.161.92.0/24
10.161.93.0/24
10.161.94.0/24
10.161.95.0/24
10.161.96.0/24
10.161.97.0/24
10.161.98.0/24
10.161.99.0/24
10.162.0.0/24
10.162.1.0/24
10.162.2.0/24
10.162.3.0/24
10.162.4.0/24
10.162.5.0/24
10.162.6.0/24
10.162.7.0/24
10.162.8.0/24
10.162.9.0/24
10.162.10.0/24
10.162.11.0/24
10.162.12.0/24
10.162.13.0/24
10.162.14.0/24
10.162.15.0/24
10.162.16.0/24
10.162.17.0/24
10.162.18.0/24
10.162.19.0/24
10.162.20.0/24
10.162.21.0/24
10.162.22.0/24
10.162.23.0/24
10.162.24.0/24
10.162.25.0/24
10.162.26.0/24
10.162.27.0/24
10.162.28.0/24
10.162.29.0/24
10.162.30.0/24
10.162.31.0/24
10.162.32.0/24
10.162.33.0/24
10.162.34.0/24
10.162.35.0/24
10.162.36.0/24
10.162.37.0/24
10.162.38.0/24
10.162.39.0/24
10.162.40.0/24
10.162.41.0/24
10.162.42.0/24
10.162.43.0/24
10.162.44.0/24
10.162.45.0/24
10.162.46.0/24
10.162.47.0/24
10.162.48.0/24
10.162.49.0/24
10.162.50.0/24
10.162.51.0/24
10.162.52.0/24
10.162.53.0/24
10.162.54.0/24
10.162.55.0/24
10.162.56.0/24
10.162.57.0/24
10.162.58.0/24
10.162.59.0/24
10.162.60.0/24
10.162.61.0/24
10.162.62.0/24
10.162.63.0/24
10.162.64.0/24
10.162.65.0/24
10.162.66.0/24
10.162.67.0/24
10.162.68.0/24
10.162.69.0/24
10.162.70.0/24
10.162.71.0/24
10.162.72.0/24
10.162.73.0/24
10.162.74.0/24
10.162.75.0/24
10.162.76.0/24
10.162.77.0/24
10.162.78.0/24
10.162.79.0/24
10.162.80.0/24
10.162.81.0/24
10.162.82.0/24
10.162.83.0/24
10.162.84.0/24
10.162.85.0/24
10.162.86.0/24
10.162.87.0/24
10.162.88.0/24
10.162.89.0/24
10.162.90.0/24
10.162.91.0/24
10.162.92.0/24
10.162.93.0/24
10.162.94.0/24
10.162.95.0/24
10.162.96.0/24
10.162.97.0/24
10.162.98.0/24
10.162.99.0/24
10.163.0.0/24
10.163.1.0/24
10.163.2.0/24
10.163.3.0/24
10.163.4.0/24
10.163.5.0/24
10.163.6.0/24
10.163.7.0/24
10.163.8.0/24
10.163.9.0/24
10.163.10.0/24
10.163.11.0/24
10.163.12.0/24
10.163.13.0/24
10.163.14.0/24
10.163.15.0/24
10.163.16.0/24
10.163.17.0/24
10.163.18.0/24
10.163.19.0/24
10.163.20.0/24
10.163.21.0/24
10.163.22.0/24
10.163.23.0/24
10.163.24.0/24
10.163.25.0/24
10.163.26.0/24
10.163.27.0/24
10.163.28.0/24
10.163.29.0/24
10.163.30.0/24
10.163.31.0/24
10.163.32.0/24
10.163.33.0/24
10.163.34.0/24
10.163.35.0/24
10.163.36.0/24
10.163.37.0/24
10.163.38.0/24
10.163.39.0/24
10.163.40.0/24
10.163.41.0/24
10.163.42.0/24
10.163.43.0/24
10.163.44.0/24
10.163.45.0/24
10.163.46.0/24
10.163.47.0/24
10.163.48.0/24
10.163.49.0/24
10.163.50.0/24
10.163.51.0/24
10.163.52.0/24
10.163.53.0/24
10.163.54.0/24
10.163.55.0/24
10.163.56.0/24
10.163.57.0/24
10.163.58.0/24
10.163.59.0/24
10.163.60.0/24
10.163.61.0/24
10.163.62.0/24
10.163.63.0/24
10.163.64.0/24
10.163.65.0/24
10.163.66.0/24
10.163.67.0/24
10.163.68.0/24
10.163.69.0/24
10.163.70.0/24
10.163.71.0/24
10.163.72.0/24
10.163.73.0/24
10.163.74.0/24
10.163.75.0/24
10.163.76.0/24
10.163.77.0/24
10.163.78.0/24
10.163.79.0/24
10.163.80.0/24
10.163.81.0/24
10.163.82.0/24
10.163.83.0/24
10.163.84.0/24
10.163.85.0/24
10.163.86.0/24
10.163.87.0/24
10.163.88.0/24
10.163.89.0/24
10.163.90.0/24
10.163.91.0/24
10.163.92.0/24
10.163.93.0/24
10.163.94.0/24
10.163.95.0/24
10.163.96.0/24
10.163.97.0/24
10.163.98.0/24
10.163.99.0/24
10.164.0.0/24
10.164.1.0/24
10.164.2.0/24
10.164.3.0/24
10.164.4.0/24
10.164.5.0/24
10.164.6.0/24
10.164.7.0/24
10.164.8.0/24
10.164.9.0/24
10.164.10.0/24
10.164.11.0/24
10.164.12.0/24
10.164.13.0/24
10.164.14.0/24
10.164.15.0/24
10.164.16.0/24
10.164.17.0/24
10.164.18.0/24
10.164.19.0/24
10.164.20.0/24
10.164.21.0/24
10.164.22.0/24
10.164.23.0/24
10.164.24.0/24
10.164.25.0/24
10.164.26.0/24
10.164.27.0/24
10.164.28.0/24
10.164.29.0/24
10.164.30.0/24
10.164.31.0/24
10.164.32.0/24
10.164.33.0/24
10.164.34.0/24
10.164.35.0/24
10.164.36.0/24
10.164.37.0/24
10.164.38.0/24
10.164.39.0/24
10.164.40.0/24
10.164.41.0/24
10.164.42.0/24
10.164.43.0/24
10.164.44.0/24
10.164.45.0/24
10.164.46.0/24
10.164.47.0/24
10.164.48.0/24
10.164.49.0/24
10.164.50.0/24
10.164.51.0/24
10.164.52.0/24
10.164.53.0/24
10.164.54.0/24
10.164.55.0/24
10.164.56.0/24
10.164.57.0/24
10.164.58.0/24
10.164.59.0/24
10.164.60.0/24
10.164.61.0/24
10.164.62.0/24
10.164.63.0/24
10.164.64.0/24
10.164.65.0/24
10.164.66.0/24
10.164.67.0/24
10.164.68.0/24
10.164.69.0/24
10.164.70.0/24
10.164.71.0/24
10.164.72.0/24
10.164.73.0/24
10.164.74.0/24
10.164.75.0/24
10.164.76.0/24
10.164.77.0/24
10.164.78.0/24
10.164.79.0/24
10.164.80.0/24
10.164.81.0/24
10.164.82.0/24
10.164.83.0/24
10.164.84.0/24
10.164.85.0/24
10.164.86.0/24
10.164.87.0/24
10.164.88.0/24
10.164.89.0/24
10.164.90.0/24
10.164.91.0/24
10.164.92.0/24
10.164.93.0/24
10.164.94.0/24
10.164.95.0/24
10.164.96.0/24
10.164.97.0/24
10.164.98.0/24
10.164.99.0/24
10.165.0.0/24
10.165.1.0/24
10.165.2.0/24
10.165.3.0/24
10.165.4.0/24
10.165.5.0/24
10.165.6.0/24
10.165.7.0/24
10.165.8.0/24
10.165.9.0/24
10.165.10.0/24
10.165.11.0/24
10.165.12.0/24
10.165.13.0/24
10.165.14.0/24
10.165.15.0/24
10.165.16.0/24
10.165.17.0/24
10.165.18.0/24
10.165.19.0/24
10.165.20.0/24
10.165.21.0/24
10.165.22.0/24
10.165.23.0/24
10.165.24.0/24
10.165.25.0/24
10.165.26.0/24
10.165.27.0/24
10.165.28.0/24
10.165.29.0/24
10.165.30.0/24
10.165.31.0/24
10.165.32.0/24
10.165.33.0/24
10.165.34.0/24
10.165.35.0/24
10.165.36.0/24
10.165.37.0/24
10.165.38.0/24
10.165.39.0/24
10.165.40.0/24
10.165.41.0/24
10.165.42.0/24
10.165.43.0/24
10.165.44.0/24
10.165.45.0/24
10.165.46.0/24
10.165.47.0/24
10.165.48.0/24
10.165.49.0/24
10.165.50.0/24
10.165.51.0/24
10.165.52.0/24
10.165.53.0/24
10.165.54.0/24
10.165.55.0/24
10.165.56.0/24
10.165.57.0/24
10.165.58.0/24
10.165.59.0/24
10.165.60.0/24
10.165.61.0/24
10.165.62.0/24
10.165.63.0/24
10.165.64.0/24
10.165.65.0/24
10.165.66.0/24
10.165.67.0/24
10.165.68.0/24
10.165.69.0/24
10.165.70.0/24
10.165.71.0/24
10.165.72.0/24
10.165.73.0/24
10.165.74.0/24
10.165.75.0/24
10.165.76.0/24`