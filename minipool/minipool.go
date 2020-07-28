package minipool

import (
    "fmt"
    "math/big"
    "sync"

    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "golang.org/x/sync/errgroup"

    "github.com/rocket-pool/rocketpool-go/rocketpool"
    rptypes "github.com/rocket-pool/rocketpool-go/types"
    "github.com/rocket-pool/rocketpool-go/utils/contract"
    "github.com/rocket-pool/rocketpool-go/utils/eth"
)


// Minipool details
type MinipoolDetails struct {
    Address common.Address              `json:"address"`
    Exists bool                         `json:"exists"`
    Pubkey rptypes.ValidatorPubkey      `json:"pubkey"`
    WithdrawalTotalBalance *big.Int     `json:"withdrawalTotalBalance"`
    WithdrawalNodeBalance *big.Int      `json:"withdrawalNodeBalance"`
    Withdrawable bool                   `json:"withdrawable"`
    WithdrawalProcessed bool            `json:"withdrawalProcessed"`
}


// Get all minipool details
func GetMinipools(rp *rocketpool.RocketPool) ([]MinipoolDetails, error) {

    // Get minipool addresses
    minipoolAddresses, err := GetMinipoolAddresses(rp)
    if err != nil {
        return []MinipoolDetails{}, err
    }

    // Data
    var wg errgroup.Group
    details := make([]MinipoolDetails, len(minipoolAddresses))

    // Load details
    for mi, minipoolAddress := range minipoolAddresses {
        mi, minipoolAddress := mi, minipoolAddress
        wg.Go(func() error {
            minipoolDetails, err := GetMinipoolDetails(rp, minipoolAddress)
            if err == nil { details[mi] = minipoolDetails }
            return err
        })
    }

    // Wait for data
    if err := wg.Wait(); err != nil {
        return []MinipoolDetails{}, err
    }

    // Return
    return details, nil

}


// Get all minipool addresses
func GetMinipoolAddresses(rp *rocketpool.RocketPool) ([]common.Address, error) {

    // Get minipool count
    minipoolCount, err := GetMinipoolCount(rp)
    if err != nil {
        return []common.Address{}, err
    }

    // Data
    var wg errgroup.Group
    addresses := make([]common.Address, minipoolCount)

    // Load addresses
    for mi := int64(0); mi < minipoolCount; mi++ {
        mi := mi
        wg.Go(func() error {
            address, err := GetMinipoolAt(rp, mi)
            if err == nil { addresses[mi] = address }
            return err
        })
    }

    // Wait for data
    if err := wg.Wait(); err != nil {
        return []common.Address{}, err
    }

    // Return
    return addresses, nil

}


// Get a node's minipool details
func GetNodeMinipools(rp *rocketpool.RocketPool, nodeAddress common.Address) ([]MinipoolDetails, error) {

    // Get minipool addresses
    minipoolAddresses, err := GetNodeMinipoolAddresses(rp, nodeAddress)
    if err != nil {
        return []MinipoolDetails{}, err
    }

    // Data
    var wg errgroup.Group
    details := make([]MinipoolDetails, len(minipoolAddresses))

    // Load details
    for mi, minipoolAddress := range minipoolAddresses {
        mi, minipoolAddress := mi, minipoolAddress
        wg.Go(func() error {
            minipoolDetails, err := GetMinipoolDetails(rp, minipoolAddress)
            if err == nil { details[mi] = minipoolDetails }
            return err
        })
    }

    // Wait for data
    if err := wg.Wait(); err != nil {
        return []MinipoolDetails{}, err
    }

    // Return
    return details, nil

}


// Get a node's minipool addresses
func GetNodeMinipoolAddresses(rp *rocketpool.RocketPool, nodeAddress common.Address) ([]common.Address, error) {

    // Get minipool count
    minipoolCount, err := GetNodeMinipoolCount(rp, nodeAddress)
    if err != nil {
        return []common.Address{}, err
    }

    // Data
    var wg errgroup.Group
    addresses := make([]common.Address, minipoolCount)

    // Load addresses
    for mi := int64(0); mi < minipoolCount; mi++ {
        mi := mi
        wg.Go(func() error {
            address, err := GetNodeMinipoolAt(rp, nodeAddress, mi)
            if err == nil { addresses[mi] = address }
            return err
        })
    }

    // Wait for data
    if err := wg.Wait(); err != nil {
        return []common.Address{}, err
    }

    // Return
    return addresses, nil

}


// Get a minipool's details
func GetMinipoolDetails(rp *rocketpool.RocketPool, minipoolAddress common.Address) (MinipoolDetails, error) {

    // Data
    var wg errgroup.Group
    var exists bool
    var pubkey rptypes.ValidatorPubkey
    var withdrawalTotalBalance *big.Int
    var withdrawalNodeBalance *big.Int
    var withdrawable bool
    var withdrawalProcessed bool

    // Load data
    wg.Go(func() error {
        var err error
        exists, err = GetMinipoolExists(rp, minipoolAddress)
        return err
    })
    wg.Go(func() error {
        var err error
        pubkey, err = GetMinipoolPubkey(rp, minipoolAddress)
        return err
    })
    wg.Go(func() error {
        var err error
        withdrawalTotalBalance, err = GetMinipoolWithdrawalTotalBalance(rp, minipoolAddress)
        return err
    })
    wg.Go(func() error {
        var err error
        withdrawalNodeBalance, err = GetMinipoolWithdrawalNodeBalance(rp, minipoolAddress)
        return err
    })
    wg.Go(func() error {
        var err error
        withdrawable, err = GetMinipoolWithdrawable(rp, minipoolAddress)
        return err
    })
    wg.Go(func() error {
        var err error
        withdrawalProcessed, err = GetMinipoolWithdrawalProcessed(rp, minipoolAddress)
        return err
    })

    // Wait for data
    if err := wg.Wait(); err != nil {
        return MinipoolDetails{}, err
    }

    // Return
    return MinipoolDetails{
        Address: minipoolAddress,
        Exists: exists,
        Pubkey: pubkey,
        WithdrawalTotalBalance: withdrawalTotalBalance,
        WithdrawalNodeBalance: withdrawalNodeBalance,
        Withdrawable: withdrawable,
        WithdrawalProcessed: withdrawalProcessed,
    }, nil

}


// Get the minipool count
func GetMinipoolCount(rp *rocketpool.RocketPool) (int64, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return 0, err
    }
    minipoolCount := new(*big.Int)
    if err := rocketMinipoolManager.Call(nil, minipoolCount, "getMinipoolCount"); err != nil {
        return 0, fmt.Errorf("Could not get minipool count: %w", err)
    }
    return (*minipoolCount).Int64(), nil
}


// Get a minipool address by index
func GetMinipoolAt(rp *rocketpool.RocketPool, index int64) (common.Address, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return common.Address{}, err
    }
    minipoolAddress := new(common.Address)
    if err := rocketMinipoolManager.Call(nil, minipoolAddress, "getMinipoolAt", big.NewInt(index)); err != nil {
        return common.Address{}, fmt.Errorf("Could not get minipool %d address: %w", index, err)
    }
    return *minipoolAddress, nil
}


// Get a node's minipool count
func GetNodeMinipoolCount(rp *rocketpool.RocketPool, nodeAddress common.Address) (int64, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return 0, err
    }
    minipoolCount := new(*big.Int)
    if err := rocketMinipoolManager.Call(nil, minipoolCount, "getNodeMinipoolCount", nodeAddress); err != nil {
        return 0, fmt.Errorf("Could not get node %s minipool count: %w", nodeAddress.Hex(), err)
    }
    return (*minipoolCount).Int64(), nil
}


// Get a node's minipool address by index
func GetNodeMinipoolAt(rp *rocketpool.RocketPool, nodeAddress common.Address, index int64) (common.Address, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return common.Address{}, err
    }
    minipoolAddress := new(common.Address)
    if err := rocketMinipoolManager.Call(nil, minipoolAddress, "getNodeMinipoolAt", nodeAddress, big.NewInt(index)); err != nil {
        return common.Address{}, fmt.Errorf("Could not get node %s minipool %d address: %w", nodeAddress.Hex(), index, err)
    }
    return *minipoolAddress, nil
}


// Get a minipool address by validator pubkey
func GetMinipoolByPubkey(rp *rocketpool.RocketPool, pubkey rptypes.ValidatorPubkey) (common.Address, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return common.Address{}, err
    }
    minipoolAddress := new(common.Address)
    if err := rocketMinipoolManager.Call(nil, minipoolAddress, "getMinipoolByPubkey", pubkey); err != nil {
        return common.Address{}, fmt.Errorf("Could not get validator %s minipool address: %w", pubkey.Hex(), err)
    }
    return *minipoolAddress, nil
}


// Check whether a minipool exists
func GetMinipoolExists(rp *rocketpool.RocketPool, minipoolAddress common.Address) (bool, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return false, err
    }
    exists := new(bool)
    if err := rocketMinipoolManager.Call(nil, exists, "getMinipoolExists", minipoolAddress); err != nil {
        return false, fmt.Errorf("Could not get minipool %s exists status: %w", minipoolAddress.Hex(), err)
    }
    return *exists, nil
}


// Get a minipool's validator pubkey
func GetMinipoolPubkey(rp *rocketpool.RocketPool, minipoolAddress common.Address) (rptypes.ValidatorPubkey, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return rptypes.ValidatorPubkey{}, err
    }
    pubkey := new(rptypes.ValidatorPubkey)
    if err := rocketMinipoolManager.Call(nil, pubkey, "getMinipoolPubkey", minipoolAddress); err != nil {
        return rptypes.ValidatorPubkey{}, fmt.Errorf("Could not get minipool %s pubkey: %w", minipoolAddress.Hex(), err)
    }
    return *pubkey, nil
}


// Get a minipool's total balance at withdrawal
func GetMinipoolWithdrawalTotalBalance(rp *rocketpool.RocketPool, minipoolAddress common.Address) (*big.Int, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return nil, err
    }
    balance := new(*big.Int)
    if err := rocketMinipoolManager.Call(nil, balance, "getMinipoolWithdrawalTotalBalance", minipoolAddress); err != nil {
        return nil, fmt.Errorf("Could not get minipool %s withdrawal total balance: %w", minipoolAddress.Hex(), err)
    }
    return *balance, nil
}


// Get a minipool's node balance at withdrawal
func GetMinipoolWithdrawalNodeBalance(rp *rocketpool.RocketPool, minipoolAddress common.Address) (*big.Int, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return nil, err
    }
    balance := new(*big.Int)
    if err := rocketMinipoolManager.Call(nil, balance, "getMinipoolWithdrawalNodeBalance", minipoolAddress); err != nil {
        return nil, fmt.Errorf("Could not get minipool %s withdrawal node balance: %w", minipoolAddress.Hex(), err)
    }
    return *balance, nil
}


// Check whether a minipool is withdrawable
func GetMinipoolWithdrawable(rp *rocketpool.RocketPool, minipoolAddress common.Address) (bool, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return false, err
    }
    withdrawable := new(bool)
    if err := rocketMinipoolManager.Call(nil, withdrawable, "getMinipoolWithdrawable", minipoolAddress); err != nil {
        return false, fmt.Errorf("Could not get minipool %s withdrawable status: %w", minipoolAddress.Hex(), err)
    }
    return *withdrawable, nil
}


// Check whether a minipool's validator withdrawal has been processed
func GetMinipoolWithdrawalProcessed(rp *rocketpool.RocketPool, minipoolAddress common.Address) (bool, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return false, err
    }
    processed := new(bool)
    if err := rocketMinipoolManager.Call(nil, processed, "getMinipoolWithdrawalProcessed", minipoolAddress); err != nil {
        return false, fmt.Errorf("Could not get minipool %s withdrawal processed status: %w", minipoolAddress.Hex(), err)
    }
    return *processed, nil
}


// Get the total length of the minipool queue
func GetQueueTotalLength(rp *rocketpool.RocketPool) (int64, error) {
    rocketMinipoolQueue, err := getRocketMinipoolQueue(rp)
    if err != nil {
        return 0, err
    }
    length := new(*big.Int)
    if err := rocketMinipoolQueue.Call(nil, length, "getTotalLength"); err != nil {
        return 0, fmt.Errorf("Could not get minipool queue total length: %w", err)
    }
    return (*length).Int64(), nil
}


// Get the total capacity of the minipool queue
func GetQueueTotalCapacity(rp *rocketpool.RocketPool) (*big.Int, error) {
    rocketMinipoolQueue, err := getRocketMinipoolQueue(rp)
    if err != nil {
        return nil, err
    }
    capacity := new(*big.Int)
    if err := rocketMinipoolQueue.Call(nil, capacity, "getTotalCapacity"); err != nil {
        return nil, fmt.Errorf("Could not get minipool queue total capacity: %w", err)
    }
    return *capacity, nil
}


// Get the capacity of the next minipool in the queue
func GetQueueNextCapacity(rp *rocketpool.RocketPool) (*big.Int, error) {
    rocketMinipoolQueue, err := getRocketMinipoolQueue(rp)
    if err != nil {
        return nil, err
    }
    capacity := new(*big.Int)
    if err := rocketMinipoolQueue.Call(nil, capacity, "getNextCapacity"); err != nil {
        return nil, fmt.Errorf("Could not get minipool queue next item capacity: %w", err)
    }
    return *capacity, nil
}


// Get the node reward amount for a minipool by node fee, user deposit balance, and staking start & end balances
func GetMinipoolNodeRewardAmount(rp *rocketpool.RocketPool, nodeFee float64, userDepositBalance, startBalance, endBalance *big.Int) (*big.Int, error) {
    rocketMinipoolStatus, err := getRocketMinipoolStatus(rp)
    if err != nil {
        return nil, err
    }
    nodeAmount := new(*big.Int)
    if err := rocketMinipoolStatus.Call(nil, nodeAmount, "getMinipoolNodeRewardAmount", eth.EthToWei(nodeFee), userDepositBalance, startBalance, endBalance); err != nil {
        return nil, fmt.Errorf("Could not get minipool node reward amount: %w", err)
    }
    return *nodeAmount, nil
}


// Submit a minipool withdrawable event
func SubmitMinipoolWithdrawable(rp *rocketpool.RocketPool, minipoolAddress common.Address, stakingStartBalance, stakingEndBalance *big.Int, opts *bind.TransactOpts) (*types.Receipt, error) {
    rocketMinipoolStatus, err := getRocketMinipoolStatus(rp)
    if err != nil {
        return nil, err
    }
    txReceipt, err := contract.Transact(rp.Client, rocketMinipoolStatus, opts, "submitMinipoolWithdrawable", minipoolAddress, stakingStartBalance, stakingEndBalance)
    if err != nil {
        return nil, fmt.Errorf("Could not submit minipool withdrawable event: %w", err)
    }
    return txReceipt, nil
}


// Get contracts
var rocketMinipoolManagerLock sync.Mutex
func getRocketMinipoolManager(rp *rocketpool.RocketPool) (*bind.BoundContract, error) {
    rocketMinipoolManagerLock.Lock()
    defer rocketMinipoolManagerLock.Unlock()
    return rp.GetContract("rocketMinipoolManager")
}
var rocketMinipoolQueueLock sync.Mutex
func getRocketMinipoolQueue(rp *rocketpool.RocketPool) (*bind.BoundContract, error) {
    rocketMinipoolQueueLock.Lock()
    defer rocketMinipoolQueueLock.Unlock()
    return rp.GetContract("rocketMinipoolQueue")
}
var rocketMinipoolStatusLock sync.Mutex
func getRocketMinipoolStatus(rp *rocketpool.RocketPool) (*bind.BoundContract, error) {
    rocketMinipoolStatusLock.Lock()
    defer rocketMinipoolStatusLock.Unlock()
    return rp.GetContract("rocketMinipoolStatus")
}

