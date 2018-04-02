package handler

import (
	"errors"
	"fmt"
	"github.com/docker/go-plugins-helpers/ipam"
	"log"
	"net"
)

type IPAMHandler struct {
	Pool map[string][]string
}

func NewHandler() IPAMHandler {
	return IPAMHandler{
		Pool: make(map[string][]string),
	}
}

// GetCapabilities returns whether or not this IPAM required pre-made MAC
func (h IPAMHandler) GetCapabilities() (response *ipam.CapabilitiesResponse, err error) {
	log.Printf("GetCapabilities called")

	return &ipam.CapabilitiesResponse{RequiresMACAddress: true}, nil
}

// GetDefaultAddressSpaces returns the default local and global address space names for this IPAM
func (h IPAMHandler) GetDefaultAddressSpaces() (response *ipam.AddressSpacesResponse, err error) {
	log.Printf("GetDefaultAddressSpaces called")
	return &ipam.AddressSpacesResponse{
		LocalDefaultAddressSpace:  "Local",
		GlobalDefaultAddressSpace: "Global",
	}, nil
}

// RequestPool handles requests for a new IP Pool
func (h IPAMHandler) RequestPool(request *ipam.RequestPoolRequest) (*ipam.RequestPoolResponse, error) {
	log.Printf("RequestPool Called: %+v", request)

	var networkName string
	if opt, ok := request.Options["network-name"]; ok {
		networkName = opt
	} else {
		return nil, errors.New("network-name is required")
	}
	if request.Pool == "" {
		return nil, errors.New("Pool is required")
	}

	ipAddr, ipNet, err := net.ParseCIDR(request.Pool)
	if err != nil {
		return nil, fmt.Errorf("Pool is invalid: %s", request.Pool)
	}
	log.Printf("Pool: %s %s", ipAddr, ipNet)

	h.Pool[networkName] = append(h.Pool[networkName], ipAddr.String())

	response := ipam.RequestPoolResponse{
		Pool: networkName,
	}
	log.Printf("RequestPoolResponse: %+v", response)

	return &response, nil
}

// ReleasePool handles requests to release an IP Pool
func (h IPAMHandler) ReleasePool(request *ipam.ReleasePoolRequest) error {
	log.Printf("ReleasePool Called: %+v", request)

	return nil
}

// RequestAddress handles requests for a new IP Address
func (h IPAMHandler) RequestAddress(request *ipam.RequestAddressRequest) (*ipam.RequestAddressResponse, error) {
	log.Printf("RequestAddress Called : %+v", request)

	addr := h.Pool[request.PoolID][0]
	response := ipam.RequestAddressResponse{
		Address: addr,
		Data:    map[string]string{},
	}

	log.Printf("RequestAddress returning %s", addr)
	return &response, nil
}

// ReleaseAddress handles requests to release an IP Address
func (h IPAMHandler) ReleaseAddress(request *ipam.ReleaseAddressRequest) (err error) {
	log.Printf("ReleaseAddress Called: %v", request)
	h.Pool[request.PoolID] = append(h.Pool[request.PoolID], request.Address)
	return nil
}
