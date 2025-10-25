package balance

type BalanceManager struct {
	repo Repository
}

type Repository interface {
	GetBalance(login string) (current, withdrawn float64, err error)
}

func New(repo Repository) *BalanceManager {
	return &BalanceManager{repo: repo}
}

func (bm *BalanceManager) GetBalance(login string) (current, withdrawn float64, err error) {
	return bm.repo.GetBalance(login)
}
