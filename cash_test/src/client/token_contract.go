package client

type TContract struct {
	Name string
}

func (d *TContract) Deposit(int) {
}

func (d *TContract) Register() {
}

func (d *TContract) BalanceOf() int {
	return 0
}

func GetTokenContract(name string) TokenContract {
	return &TContract{name}
}
