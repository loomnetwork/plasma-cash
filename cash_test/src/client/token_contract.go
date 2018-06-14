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

func (d *TContract) Account() *Account {
	return &Account{}
}

func GetTokenContract(name string) TokenContract {
	return &TContract{name}
}
