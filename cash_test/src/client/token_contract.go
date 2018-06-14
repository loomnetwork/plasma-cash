package client

type DummyContract struct {
	Name string
}

func (d *DummyContract) Deposit(int) {
}

func (d *DummyContract) Register() {
}

func GetTokenContract(name string) TokenContract {
	return &DummyContract{name}
}
