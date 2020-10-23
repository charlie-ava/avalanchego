// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avm

import (
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/api"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow/choices"
	"github.com/ava-labs/avalanchego/utils/formatting"
	cjson "github.com/ava-labs/avalanchego/utils/json"
	"github.com/ava-labs/avalanchego/utils/rpc"
)

// Client ...
type Client struct {
	requester rpc.EndpointRequester
}

// NewClient returns an AVM client for interacting with avm [chain]
func NewClient(uri, chain string, requestTimeout time.Duration) *Client {
	return &Client{
		requester: rpc.NewEndpointRequester(uri, fmt.Sprintf("/ext/bc/%s", chain), "avm", requestTimeout),
	}
}

// IssueTx issues a transaction to a node and returns the TxID
func (c *Client) IssueTx(txBytes []byte) (ids.ID, error) {
	res := &api.JSONTxID{}
	err := c.requester.SendRequest("issueTx", &api.FormattedTx{
		Tx:       formatting.Hex{Bytes: txBytes}.String(),
		Encoding: formatting.HexEncoding,
	}, res)
	if err != nil {
		return ids.Empty, err
	}
	return res.TxID, nil
}

func (c *Client) GetTxStatus(txID ids.ID) (choices.Status, error) {
	res := &GetTxStatusReply{}
	err := c.requester.SendRequest("getTxStatus", &api.JSONTxID{
		TxID: txID,
	}, res)
	if err != nil {
		return choices.Unknown, err
	}
	return res.Status, nil
}

func (c *Client) GetTx(txID ids.ID) ([]byte, error) {
	res := &api.FormattedTx{}
	err := c.requester.SendRequest("getTx", &api.GetTxArgs{
		TxID:     txID,
		Encoding: formatting.HexEncoding,
	}, res)
	if err != nil {
		return nil, err
	}

	formatter := formatting.Hex{}
	if err := formatter.FromString(res.Tx); err != nil {
		return nil, err
	}
	return formatter.Bytes, nil
}

// GetUTXOs returns the byte representation of the UTXOs controlled by [addrs]
func (c *Client) GetUTXOs(addrs []string, limit uint32, startAddress, startUTXOID string) ([][]byte, Index, error) {
	res := &GetUTXOsReply{}
	err := c.requester.SendRequest("getUTXOs", &GetUTXOsArgs{
		Addresses: addrs,
		Limit:     cjson.Uint32(limit),
		StartIndex: Index{
			Address: startAddress,
			UTXO:    startUTXOID,
		},
		Encoding: formatting.HexEncoding,
	}, res)
	if err != nil {
		return nil, Index{}, err
	}

	formatter := formatting.Hex{}
	utxos := make([][]byte, len(res.UTXOs))
	for i, utxo := range res.UTXOs {
		if err := formatter.FromString(utxo); err != nil {
			return nil, Index{}, err
		}
		utxos[i] = formatter.Bytes
	}
	return utxos, res.EndIndex, nil
}

func (c *Client) GetAssetDescription(assetID string) (*GetAssetDescriptionReply, error) {
	res := &GetAssetDescriptionReply{}
	err := c.requester.SendRequest("getAssetDescription", &GetAssetDescriptionArgs{
		AssetID: assetID,
	}, res)
	return res, err
}

func (c *Client) GetBalance(addr string, assetID string) (*GetBalanceReply, error) {
	res := &GetBalanceReply{}
	err := c.requester.SendRequest("getBalance", &GetBalanceArgs{
		Address: addr,
		AssetID: assetID,
	}, res)
	return res, err
}

func (c *Client) GetAllBalances(addr string, assetID string) (*GetAllBalancesReply, error) {
	res := &GetAllBalancesReply{}
	err := c.requester.SendRequest("getAllBalances", &api.JSONAddress{
		Address: addr,
	}, res)
	return res, err
}

func (c *Client) CreateFixedCapAsset(
	user api.UserPass,
	name,
	symbol string,
	denomination byte,
	holders []*Holder,
	from []string,
	changeAddr string,
) (ids.ID, error) {
	res := &FormattedAssetID{}
	err := c.requester.SendRequest("createFixedCapAsset", &CreateAssetArgs{
		JSONSpendHeader: api.JSONSpendHeader{
			UserPass:       user,
			JSONFromAddrs:  api.JSONFromAddrs{From: from},
			JSONChangeAddr: api.JSONChangeAddr{ChangeAddr: changeAddr},
		},
		Name:           name,
		Symbol:         symbol,
		Denomination:   denomination,
		InitialHolders: holders,
	}, res)
	if err != nil {
		return ids.Empty, err
	}
	return res.AssetID, nil
}

func (c *Client) CreateVariableCapAsset(
	user api.UserPass,
	name,
	symbol string,
	denomination byte,
	minters []Owners,
	from []string,
	changeAddr string,
) (ids.ID, error) {
	res := &FormattedAssetID{}
	err := c.requester.SendRequest("createVariableCapAsset", &CreateAssetArgs{
		JSONSpendHeader: api.JSONSpendHeader{
			UserPass:       user,
			JSONFromAddrs:  api.JSONFromAddrs{From: from},
			JSONChangeAddr: api.JSONChangeAddr{ChangeAddr: changeAddr},
		},
		Name:         name,
		Symbol:       symbol,
		Denomination: denomination,
		MinterSets:   minters,
	}, res)
	if err != nil {
		return ids.Empty, err
	}
	return res.AssetID, nil
}

func (c *Client) CreateNFTAsset(
	user api.UserPass,
	name,
	symbol string,
	minters []Owners,
	from []string,
	changeAddr string,
) (ids.ID, error) {
	res := &FormattedAssetID{}
	err := c.requester.SendRequest("createNFTAsset", &CreateNFTAssetArgs{
		JSONSpendHeader: api.JSONSpendHeader{
			UserPass:       user,
			JSONFromAddrs:  api.JSONFromAddrs{From: from},
			JSONChangeAddr: api.JSONChangeAddr{ChangeAddr: changeAddr},
		},
		Name:       name,
		Symbol:     symbol,
		MinterSets: minters,
	}, res)
	if err != nil {
		return ids.Empty, err
	}
	return res.AssetID, nil
}

func (c *Client) CreateAddress(user api.UserPass) (string, error) {
	res := &api.JSONAddress{}
	err := c.requester.SendRequest("createAddress", &user, res)
	if err != nil {
		return "", err
	}
	return res.Address, nil
}

func (c *Client) ListAddresses(user api.UserPass) ([]string, error) {
	res := &api.JSONAddresses{}
	err := c.requester.SendRequest("listAddresses", &user, res)
	if err != nil {
		return nil, err
	}
	return res.Addresses, nil
}

func (c *Client) ExportKey(user api.UserPass, addr string) (string, error) {
	res := &ExportKeyReply{}
	err := c.requester.SendRequest("exportKey", &ExportKeyArgs{
		UserPass: user,
		Address:  addr,
	}, res)
	if err != nil {
		return "", err
	}
	return res.PrivateKey, nil
}

func (c *Client) ImportKey(user api.UserPass, privateKey string) (string, error) {
	res := &api.JSONAddress{}
	err := c.requester.SendRequest("importKey", &ImportKeyArgs{
		UserPass:   user,
		PrivateKey: privateKey,
	}, res)
	if err != nil {
		return "", err
	}
	return res.Address, nil
}

func (c *Client) Send(
	user api.UserPass,
	amount uint64,
	assetID,
	to string,
	from []string,
	changeAddr string,
) (ids.ID, error) {
	res := &api.JSONTxID{}
	err := c.requester.SendRequest("send", &SendArgs{
		JSONSpendHeader: api.JSONSpendHeader{
			UserPass:       user,
			JSONFromAddrs:  api.JSONFromAddrs{From: from},
			JSONChangeAddr: api.JSONChangeAddr{ChangeAddr: changeAddr},
		},
		SendOutput: SendOutput{
			Amount:  cjson.Uint64(amount),
			AssetID: assetID,
			To:      to,
		},
	}, res)
	if err != nil {
		return ids.Empty, err
	}
	return res.TxID, nil
}

func (c *Client) SendMultiple(
	user api.UserPass,
	outputs []SendOutput,
	from []string,
	changeAddr string,
) (ids.ID, error) {
	res := &api.JSONTxID{}
	err := c.requester.SendRequest("send", &SendMultipleArgs{
		JSONSpendHeader: api.JSONSpendHeader{
			UserPass:       user,
			JSONFromAddrs:  api.JSONFromAddrs{From: from},
			JSONChangeAddr: api.JSONChangeAddr{ChangeAddr: changeAddr},
		},
		Outputs: outputs,
	}, res)
	if err != nil {
		return ids.Empty, err
	}
	return res.TxID, nil
}

func (c *Client) Mint(
	user api.UserPass,
	amount uint64,
	assetID,
	to string,
	from []string,
	changeAddr string,
) (ids.ID, error) {
	res := &api.JSONTxID{}
	err := c.requester.SendRequest("mint", &MintArgs{
		JSONSpendHeader: api.JSONSpendHeader{
			UserPass:       user,
			JSONFromAddrs:  api.JSONFromAddrs{From: from},
			JSONChangeAddr: api.JSONChangeAddr{ChangeAddr: changeAddr},
		},
		Amount:  cjson.Uint64(amount),
		AssetID: assetID,
		To:      to,
	}, res)
	if err != nil {
		return ids.Empty, err
	}
	return res.TxID, nil
}

func (c *Client) SendNFT(
	user api.UserPass,
	assetID string,
	groupID uint32,
	to string,
	from []string,
	changeAddr string,
) (ids.ID, error) {
	res := &api.JSONTxID{}
	err := c.requester.SendRequest("sendNFT", &SendNFTArgs{
		JSONSpendHeader: api.JSONSpendHeader{
			UserPass:       user,
			JSONFromAddrs:  api.JSONFromAddrs{From: from},
			JSONChangeAddr: api.JSONChangeAddr{ChangeAddr: changeAddr},
		},
		AssetID: assetID,
		GroupID: cjson.Uint32(groupID),
		To:      to,
	}, res)
	if err != nil {
		return ids.Empty, err
	}
	return res.TxID, nil
}

func (c *Client) MintNFT(
	user api.UserPass,
	assetID string,
	payload []byte,
	to string,
	from []string,
	changeAddr string,
) (ids.ID, error) {
	res := &api.JSONTxID{}
	err := c.requester.SendRequest("mintNFT", &MintNFTArgs{
		JSONSpendHeader: api.JSONSpendHeader{
			UserPass:       user,
			JSONFromAddrs:  api.JSONFromAddrs{From: from},
			JSONChangeAddr: api.JSONChangeAddr{ChangeAddr: changeAddr},
		},
		AssetID:  assetID,
		Payload:  formatting.Hex{Bytes: payload}.String(),
		Encoding: formatting.HexEncoding,
		To:       to,
	}, res)
	if err != nil {
		return ids.Empty, err
	}
	return res.TxID, nil
}

func (c *Client) ImportAVAX(user api.UserPass, to, sourceChain string) (ids.ID, error) {
	res := &api.JSONTxID{}
	err := c.requester.SendRequest("importAVAX", &ImportArgs{
		UserPass:    user,
		To:          to,
		SourceChain: sourceChain,
	}, res)
	if err != nil {
		return ids.Empty, err
	}
	return res.TxID, nil
}

func (c *Client) ExportAVAX(
	user api.UserPass,
	amount uint64,
	to string,
	from []string,
	changeAddr string,
) (ids.ID, error) {
	res := &api.JSONTxID{}
	err := c.requester.SendRequest("exportAVAX", &ExportAVAXArgs{
		JSONSpendHeader: api.JSONSpendHeader{
			UserPass:       user,
			JSONFromAddrs:  api.JSONFromAddrs{From: from},
			JSONChangeAddr: api.JSONChangeAddr{ChangeAddr: changeAddr},
		},
		Amount: cjson.Uint64(amount),
		To:     to,
	}, res)
	if err != nil {
		return ids.Empty, err
	}
	return res.TxID, nil
}

func (c *Client) Import(user api.UserPass, to, sourceChain string) (ids.ID, error) {
	res := &api.JSONTxID{}
	err := c.requester.SendRequest("importAVAX", &ImportArgs{
		UserPass:    user,
		To:          to,
		SourceChain: sourceChain,
	}, res)
	return res.TxID, err
}

func (c *Client) Export(
	user api.UserPass,
	amount uint64,
	to string,
	assetID string,
	from []string,
	changeAddr string,
) (ids.ID, error) {
	res := &api.JSONTxID{}
	err := c.requester.SendRequest("exportAVAX", &ExportArgs{
		ExportAVAXArgs: ExportAVAXArgs{
			JSONSpendHeader: api.JSONSpendHeader{
				UserPass:       user,
				JSONFromAddrs:  api.JSONFromAddrs{From: from},
				JSONChangeAddr: api.JSONChangeAddr{ChangeAddr: changeAddr},
			},
			Amount: cjson.Uint64(amount),
			To:     to,
		},
		AssetID: assetID,
	}, res)

	return res.TxID, err
}
