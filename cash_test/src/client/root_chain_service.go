package client

import (
	"crypto/ecdsa"
	"ethcontract"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RootChainService struct {
	Name           string
	plasmaContract *ethcontract.RootChain
	callerKey      *ecdsa.PrivateKey
	callerAddr     common.Address
	transactOpts   *bind.TransactOpts
	callOpts       *bind.CallOpts
}

func (d *RootChainService) PlasmaCoin(slot uint64) (*PlasmaCoin, error) {
	uid, depositBlockNum, denom, ownerAddr, state, err := d.plasmaContract.GetPlasmaCoin(&bind.CallOpts{
		From: d.callerAddr,
	}, slot)
	if err != nil {
		return nil, err
	}
	return &PlasmaCoin{
		UID:             uid,
		DepositBlockNum: depositBlockNum.Int64(),
		Denomination:    denom,
		Owner:           ownerAddr.Hex(),
		State:           PlasmaCoinState(state),
	}, nil
}

func (d *RootChainService) Withdraw(slot uint64) error {
	_, err := d.plasmaContract.Withdraw(d.transactOpts, slot)
	return err
}

func (d *RootChainService) ChallengeBefore(slot uint64, prevTx Tx, exitingTx Tx,
	prevTxInclusionProof Proof, exitingTxInclusionProof Proof,
	sig []byte, prevTxBlockNum int64, exitingTxBlockNum int64) ([]byte, error) {
	prevTxBytes, err := prevTx.RlpEncode()
	if err != nil {
		return nil, err
	}
	exitingTxBytes, err := exitingTx.RlpEncode()
	if err != nil {
		return nil, err
	}
	tx, err := d.plasmaContract.ChallengeBefore(
		d.transactOpts, slot, prevTxBytes, exitingTxBytes,
		prevTxInclusionProof, exitingTxInclusionProof, sig,
		big.NewInt(prevTxBlockNum), big.NewInt(exitingTxBlockNum))
	if err != nil {
		return nil, err
	}
	return tx.Hash().Bytes(), nil
}

func (d *RootChainService) RespondChallengeBefore(slot uint64, challengingBlockNumber int64,
	challengingTx Tx, proof Proof) ([]byte, error) {
	challengingTxBytes, err := challengingTx.RlpEncode()
	if err != nil {
		return nil, err
	}
	tx, err := d.plasmaContract.RespondChallengeBefore(
		d.transactOpts, slot, big.NewInt(challengingBlockNumber), challengingTxBytes, proof)
	if err != nil {
		return nil, err
	}
	return tx.Hash().Bytes(), nil
}

func (d *RootChainService) ChallengeBetween(slot uint64, challengingBlockNumber int64,
	challengingTx Tx, proof Proof, sig []byte) ([]byte, error) {
	challengingTxBytes, err := challengingTx.RlpEncode()
	if err != nil {
		return nil, err
	}
	tx, err := d.plasmaContract.ChallengeBetween(
		d.transactOpts, slot, big.NewInt(challengingBlockNumber), challengingTxBytes, proof, sig)
	if err != nil {
		return nil, err
	}
	return tx.Hash().Bytes(), nil
}

func (d *RootChainService) ChallengeAfter(slot uint64, challengingBlockNumber int64,
	challengingTx Tx, proof Proof, sig []byte) ([]byte, error) {
	challengingTxBytes, err := challengingTx.RlpEncode()
	if err != nil {
		return nil, err
	}
	tx, err := d.plasmaContract.ChallengeAfter(
		d.transactOpts, slot, big.NewInt(challengingBlockNumber), challengingTxBytes, proof, sig)
	if err != nil {
		return nil, err
	}
	return tx.Hash().Bytes(), nil
}

func (d *RootChainService) StartExit(
	slot uint64, prevTx Tx, exitingTx Tx, prevTxInclusion Proof, exitingTxInclusion Proof,
	sigs []byte, prevTxIncBlock int64, exitingTxIncBlock int64) ([]byte, error) {
	prevTxBytes, err := prevTx.RlpEncode()
	if err != nil {
		return nil, err
	}
	exitingTxBytes, err := exitingTx.RlpEncode()
	if err != nil {
		return nil, err
	}
	d.transactOpts.Value = big.NewInt(100000000000000000) //0.1 eth

	fmt.Printf("\nprevTxInclusion.Bytes()-%v-len(%d)\n", prevTxInclusion, len(prevTxInclusion))
	tx, err := d.plasmaContract.StartExit(
		d.transactOpts, slot,
		prevTxBytes, exitingTxBytes, prevTxInclusion, exitingTxInclusion,
		sigs, big.NewInt(prevTxIncBlock), big.NewInt(exitingTxIncBlock))

	d.transactOpts.Value = big.NewInt(0)
	if err != nil {
		return nil, err
	}
	return tx.Hash().Bytes(), nil
}

func (d *RootChainService) FinalizeExits() error {
	_, err := d.plasmaContract.FinalizeExits(d.transactOpts)
	return err
}

func (d *RootChainService) WithdrawBonds() error {
	_, err := d.plasmaContract.WithdrawBonds(d.transactOpts)
	return err
}

func (d *RootChainService) SubmitBlock(blockNum *big.Int, merkleRoot [32]byte) error {
	_, err := d.plasmaContract.SubmitBlock(d.transactOpts, merkleRoot)
	return err
}

func (d *RootChainService) DebugCoinMetaData() {
	coins, err := d.plasmaContract.NumCoins(d.callOpts) //todo make this readonly
	fmt.Printf("Num coins -%v\n", coins)
	if err != nil {
		panic(err)
	}
	for x := uint64(0); x < coins; x++ {
		//uid, c.depositBlock, c.denomination, c.owner, c.state
		uid, _, _, _, state, err := d.plasmaContract.GetPlasmaCoin(d.callOpts, x)
		fmt.Printf("Num coins -%d -(uid)-%v -(state)-%v\n", x, uid, state)

		if err != nil {
			panic(err)
		}
	}
}

var conn *ethclient.Client

func InitClients(connStr string) {
	var err error
	conn, err = ethclient.Dial(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
}

func NewRootChainService(callerName string, callerKey *ecdsa.PrivateKey, boundContract *ethcontract.RootChain) *RootChainService {
	auth := bind.NewKeyedTransactor(callerKey)
	// If gas price isn't set explicitely then go-ethereum will attempt to query the suggested gas
	// price, unfortunatley ganache-cli v6.1.2 seems to encode the gas price in a format go-ethereum
	// can't decode correctly, so this error is returned whenver you attempt to call a contract:
	// failed to suggest gas price: json: cannot unmarshal hex number with leading zero digits into Go value of type *hexutil.Big
	//
	// Earlier versions of ganache-cli don't seem to exhibit this issue, but they're broken in other
	// ways (logs aren't hex-encoded correctly).
	auth.GasPrice = big.NewInt(20000)
	auth.GasLimit = uint64(3141592)
	callerAddr := crypto.PubkeyToAddress(callerKey.PublicKey)
	return &RootChainService{
		Name:           callerName,
		callerKey:      callerKey,
		callerAddr:     callerAddr,
		plasmaContract: boundContract,
		transactOpts:   auth,
		callOpts: &bind.CallOpts{
			From: callerAddr,
		},
	}
}
