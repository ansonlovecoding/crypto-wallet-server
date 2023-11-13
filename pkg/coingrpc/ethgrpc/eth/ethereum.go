package eth

import (
	"Share-Wallet/pkg/wallet/coin"
	"Share-Wallet/pkg/wallet/config"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Ethereum includes client to call JSON-RPC
type Ethereum struct {
	ethClient    *ethclient.Client
	rpcClient    *ethrpc.Client
	erc20        *ERC20
	chainConf    *chaincfg.Params
	coinTypeCode coin.CoinTypeCode
	ctx          context.Context
	netID        uint16
	version      string
	keyDir       string
	isParity     bool
}

// NewEthereum creates ethereum object
func NewEthereum(
	ctx context.Context,
	ethClient *ethclient.Client,
	rpcClient *ethrpc.Client,
	erc20 *ERC20,
	coinTypeCode coin.CoinTypeCode,
	conf *config.Ethereum) (*Ethereum, error) {

	eth := &Ethereum{
		ethClient:    ethClient,
		rpcClient:    rpcClient,
		erc20:        erc20,
		coinTypeCode: coinTypeCode,
		ctx:          ctx,
		keyDir:       conf.KeyDirName,
	}

	// key dir
	if eth.keyDir == "" {
		//admin_datadir in some version is not available
		dirName, _ := eth.AdminDataDir()
		if dirName != "" {
			eth.keyDir = fmt.Sprintf("%s/keystore", dirName)
		}

	}
	log.Print(zap.String("eth.keyDir", eth.keyDir))

	// get NetID
	netID, err := eth.NetVersion()
	if err != nil {
		return nil, errors.Wrap(err, "fail to call eth.NetVersion()")
	}
	eth.netID = netID

	if netID == 1 {
		eth.chainConf = &chaincfg.MainNetParams
	} else {
		eth.chainConf = &chaincfg.TestNet3Params
	}

	// get client version
	clientVer, err := eth.ClientVersion()
	if err != nil {
		return nil, errors.Wrap(err, "fail to call eth.ClientVersion()")
	}
	eth.version = clientVer

	eth.isParity = isParity(clientVer)

	// check sync progress
	res, isSyncing, err := eth.Syncing()
	if err != nil {
		return nil, errors.Wrap(err, "fail to call eth.Syncing()")
	}
	if isSyncing {
		log.Print("sync is not completed yet")
	}
	if res != nil {
		log.Print("still syncing",
			zap.Int64("startingBlock", res.StartingBlock),
			zap.Int64("currentBlock", res.CurrentBlock),
			zap.Int64("highestBlock", res.HighestBlock),
		)
	}

	// check network connections
	isListening, err := eth.NetListening()
	if err != nil {
		return nil, errors.Wrap(err, "fail to call eth.NetListening()")
	}
	if !isListening {
		log.Print("network is not working")
	}

	return eth, nil
}

// Close disconnect to server
func (e *Ethereum) Close() {
	if e.rpcClient != nil {
		e.rpcClient.Close()
	}
}

// CoinTypeCode returns coinTypeCode
func (e *Ethereum) CoinTypeCode() coin.CoinTypeCode {
	return e.coinTypeCode
}

// GetChainConf returns chain conf
func (e *Ethereum) GetChainConf() *chaincfg.Params {
	return e.chainConf
}

func isParity(target string) bool {
	return strings.Contains(target, ClientVersionParity.String())
}

// GetUSDTERC20ContractAddress returns the contract address of USDT-ERC20
func (e *Ethereum) GetUSDTERC20ContractAddress() string {
	return e.erc20.contractAddress
}

func (e *Ethereum) GetUSDTERC20AbiJSON() string {
	return e.erc20.abiJSON
}
