package client

import (
	"crypto/ecdsa"
	"ethcontract"
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

	//return self.sign_and_send(self.contract.functions.challengeBefore, args,	  value=self.BOND)
	return []byte{}, nil
}

func (d *RootChainService) RespondChallengeBefore(slot uint64, challengingBlockNumber int64,
	challenging_transaction Tx, proof Proof) ([]byte, error) {

	//return self.sign_and_send(self.contract.functions.respond_challenge_before,					args)
	return []byte{}, nil
}

func (d *RootChainService) ChallengeBetween(slot uint64, challengingBlockNumber int64,
	challengingTransaction Tx, proof Proof) ([]byte, error) {
	//return self.sign_and_send(self.contract.functions.challengeBetween, args)
	return []byte{}, nil
}

func (d *RootChainService) ChallengeAfter(slot uint64, challengingBlockNumber int64,
	challengingTransaction Tx, proof Proof) ([]byte, error) {
	//return self.sign_and_send(self.contract.functions.challengeAfter, args)
	return []byte{}, nil
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
	_, err = d.plasmaContract.StartExit(
		d.transactOpts, slot,
		prevTxBytes, exitingTxBytes, prevTxInclusion.Bytes(), exitingTxInclusion.Bytes(),
		sigs, big.NewInt(prevTxIncBlock), big.NewInt(exitingTxIncBlock))
	return []byte{}, err
}

func (d *RootChainService) FinalizeExits() error {
	_, err := d.plasmaContract.FinalizeExits(d.transactOpts)
	return err
}

// TODO: implement for challenge_after_demo
func (d *RootChainService) WithdrawBonds() error {
	return nil
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
	return &RootChainService{
		Name:           callerName,
		callerKey:      callerKey,
		callerAddr:     crypto.PubkeyToAddress(callerKey.PublicKey),
		plasmaContract: boundContract,
		transactOpts:   auth,
	}
}
