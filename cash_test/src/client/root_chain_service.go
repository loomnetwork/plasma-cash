package client

type RootChainService struct {
	Name string
}

func (d *RootChainService) PlasmaCoin(int) {
}

func (d *RootChainService) Withdraw(int) {
}

func (d *RootChainService) FinalizeExits() {
}
func (d *RootChainService) WithdrawBonds() {
}

func GetRootChain(name string) RootChainClient {
	return &RootChainService{name}
}
