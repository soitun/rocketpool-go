package protocol

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

// Bootstrap a bool setting
func BootstrapBool(rp *rocketpool.RocketPool, contractName, settingPath string, value bool, opts *bind.TransactOpts) (common.Hash, error) {
    rocketDAOProtocol, err := getRocketDAOProtocol(rp)
    if err != nil {
        return common.Hash{}, err
    }
    hash, err := rocketDAOProtocol.Transact(opts, "bootstrapSettingBool", contractName, settingPath, value)
    if err != nil {
        return common.Hash{}, fmt.Errorf("Could not bootstrap protocol setting %s.%s: %w", contractName, settingPath, err)
    }
    return hash, nil
}


// Bootstrap a uint256 setting
func BootstrapUint(rp *rocketpool.RocketPool, contractName, settingPath string, value *big.Int, opts *bind.TransactOpts) (common.Hash, error) {
    rocketDAOProtocol, err := getRocketDAOProtocol(rp)
    if err != nil {
        return common.Hash{}, err
    }
    hash, err := rocketDAOProtocol.Transact(opts, "bootstrapSettingUint", contractName, settingPath, value)
    if err != nil {
        return common.Hash{}, fmt.Errorf("Could not bootstrap protocol setting %s.%s: %w", contractName, settingPath, err)
    }
    return hash, nil
}


// Bootstrap an address setting
func BootstrapAddress(rp *rocketpool.RocketPool, contractName, settingPath string, value common.Address, opts *bind.TransactOpts) (common.Hash, error) {
    rocketDAOProtocol, err := getRocketDAOProtocol(rp)
    if err != nil {
        return common.Hash{}, err
    }
    hash, err := rocketDAOProtocol.Transact(opts, "bootstrapSettingAddress", contractName, settingPath, value)
    if err != nil {
        return common.Hash{}, fmt.Errorf("Could not bootstrap protocol setting %s.%s: %w", contractName, settingPath, err)
    }
    return hash, nil
}


// Bootstrap a rewards claimer
func BootstrapClaimer(rp *rocketpool.RocketPool, contractName string, amount float64, opts *bind.TransactOpts) (common.Hash, error) {
    rocketDAOProtocol, err := getRocketDAOProtocol(rp)
    if err != nil {
        return common.Hash{}, err
    }
    hash, err := rocketDAOProtocol.Transact(opts, "bootstrapSettingClaimer", contractName, eth.EthToWei(amount))
    if err != nil {
        return common.Hash{}, fmt.Errorf("Could not bootstrap protocol rewards claimer %s: %w", contractName, err)
    }
    return hash, nil
}


// Get contracts
var rocketDAOProtocolLock sync.Mutex
func getRocketDAOProtocol(rp *rocketpool.RocketPool) (*rocketpool.Contract, error) {
    rocketDAOProtocolLock.Lock()
    defer rocketDAOProtocolLock.Unlock()
    return rp.GetContract("rocketDAOProtocol")
}

