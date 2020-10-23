package info

import (
	"time"

	"github.com/ava-labs/avalanchego/network"
	"github.com/ava-labs/avalanchego/utils/rpc"
)

// Client is an Info API Client
type Client struct {
	requester rpc.EndpointRequester
}

// NewClient returns a new Info API Client
func NewClient(uri string, requestTimeout time.Duration) *Client {
	return &Client{
		requester: rpc.NewEndpointRequester(uri, "/ext/info", "info", requestTimeout),
	}
}

// GetNodeID ...
func (c *Client) GetNodeID() (string, error) {
	res := &GetNodeIDReply{}
	err := c.requester.SendRequest("getNodeID", struct{}{}, res)
	return res.NodeID, err
}

// GetNetworkID ...
func (c *Client) GetNetworkID() (uint32, error) {
	res := &GetNetworkIDReply{}
	err := c.requester.SendRequest("getNetworkID", struct{}{}, res)
	return uint32(res.NetworkID), err
}

// GetNetworkName ...
func (c *Client) GetNetworkName() (string, error) {
	res := &GetNetworkNameReply{}
	err := c.requester.SendRequest("getNetworkName", struct{}{}, res)
	return res.NetworkName, err
}

// GetBlockchainID ...
func (c *Client) GetBlockchainID() (string, error) {
	res := &GetBlockchainIDReply{}
	err := c.requester.SendRequest("getBlockchainID", struct{}{}, res)
	return res.BlockchainID, err
}

// Peers ...
func (c *Client) Peers() ([]network.PeerID, error) {
	res := &PeersReply{}
	err := c.requester.SendRequest("peers", struct{}{}, res)
	return res.Peers, err
}

// IsBootstrapped ...
func (c *Client) IsBootstrapped(chain string) (bool, error) {
	res := &IsBootstrappedResponse{}
	err := c.requester.SendRequest("isBootstrapped", &IsBootstrappedArgs{
		Chain: chain,
	}, res)
	return res.IsBootstrapped, err
}
