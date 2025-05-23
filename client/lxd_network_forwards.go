package lxd

import (
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/canonical/lxd/shared/api"
)

// GetNetworkForwardAddresses returns a list of network forward listen addresses.
func (r *ProtocolLXD) GetNetworkForwardAddresses(networkName string) ([]string, error) {
	err := r.CheckExtension("network_forward")
	if err != nil {
		return nil, err
	}

	// Fetch the raw URL values.
	urls := []string{}
	baseURL := "/networks/" + url.PathEscape(networkName) + "/forwards"
	_, err = r.queryStruct(http.MethodGet, baseURL, nil, "", &urls)
	if err != nil {
		return nil, err
	}

	// Parse it.
	return urlsToResourceNames(baseURL, urls...)
}

// GetNetworkForwards returns a list of Network forward structs.
func (r *ProtocolLXD) GetNetworkForwards(networkName string) ([]api.NetworkForward, error) {
	err := r.CheckExtension("network_forward")
	if err != nil {
		return nil, err
	}

	forwards := []api.NetworkForward{}

	// Fetch the raw value.
	_, err = r.queryStruct(http.MethodGet, "/networks/"+url.PathEscape(networkName)+"/forwards?recursion=1", nil, "", &forwards)
	if err != nil {
		return nil, err
	}

	return forwards, nil
}

// GetNetworkForward returns a Network forward entry for the provided network and listen address.
func (r *ProtocolLXD) GetNetworkForward(networkName string, listenAddress string) (*api.NetworkForward, string, error) {
	err := r.CheckExtension("network_forward")
	if err != nil {
		return nil, "", err
	}

	forward := api.NetworkForward{}

	// Fetch the raw value.
	etag, err := r.queryStruct(http.MethodGet, "/networks/"+url.PathEscape(networkName)+"/forwards/"+url.PathEscape(listenAddress), nil, "", &forward)
	if err != nil {
		return nil, "", err
	}

	return &forward, etag, nil
}

// CreateNetworkForward defines a new network forward using the provided struct.
func (r *ProtocolLXD) CreateNetworkForward(networkName string, forward api.NetworkForwardsPost) error {
	err := r.CheckExtension("network_forward")
	if err != nil {
		return err
	}

	listenAddressIP := net.ParseIP(forward.ListenAddress)
	if listenAddressIP == nil {
		return fmt.Errorf("Invalid network forward listen address: %s", forward.ListenAddress)
	}

	if listenAddressIP.IsUnspecified() {
		err := r.CheckExtension("network_allocate_external_ips")
		if err != nil {
			return err
		}
	}

	// Send the request.
	_, _, err = r.query(http.MethodPost, "/networks/"+url.PathEscape(networkName)+"/forwards", forward, "")
	if err != nil {
		return err
	}

	return nil
}

// UpdateNetworkForward updates the network forward to match the provided struct.
func (r *ProtocolLXD) UpdateNetworkForward(networkName string, listenAddress string, forward api.NetworkForwardPut, ETag string) error {
	err := r.CheckExtension("network_forward")
	if err != nil {
		return err
	}

	// Send the request.
	_, _, err = r.query(http.MethodPut, "/networks/"+url.PathEscape(networkName)+"/forwards/"+url.PathEscape(listenAddress), forward, ETag)
	if err != nil {
		return err
	}

	return nil
}

// DeleteNetworkForward deletes an existing network forward.
func (r *ProtocolLXD) DeleteNetworkForward(networkName string, listenAddress string) error {
	err := r.CheckExtension("network_forward")
	if err != nil {
		return err
	}

	// Send the request.
	_, _, err = r.query(http.MethodDelete, "/networks/"+url.PathEscape(networkName)+"/forwards/"+url.PathEscape(listenAddress), nil, "")
	if err != nil {
		return err
	}

	return nil
}
