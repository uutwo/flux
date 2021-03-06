package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yehohanan7/flux"
	"github.com/yehohanan7/flux/cqrs"
)

type Account struct {
	cqrs.Aggregate
	Id      string
	Balance int
}

type AccountCredited struct {
	Amount int
}

type AccountDebited struct {
	Amount int
}

func (acc *Account) HandleAccountCredited(event AccountCredited) {
	acc.Balance = acc.Balance + event.Amount
}

func (acc *Account) HandleAccountDebited(event AccountDebited) {
	acc.Balance = acc.Balance - event.Amount
}

func (acc *Account) Credit(amount int) {
	acc.Update(AccountCredited{amount})
}

func (acc *Account) Debit(amount int) {
	acc.Update(AccountDebited{amount})
}

func TestAggregate(t *testing.T) {
	for _, store := range []cqrs.EventStore{flux.NewBoltStore("flux_smoke.db"), flux.NewMemoryStore()} {
		accountId := "account-id"
		existingAccount := new(Account)
		existingAccount.Aggregate = cqrs.NewAggregate(accountId, existingAccount, store)

		existingAccount.Credit(5)
		existingAccount.Credit(10)
		existingAccount.Debit(1)
		existingAccount.Save()

		account := new(Account)
		account.Aggregate = cqrs.GetAggregate(accountId, account, store)
		assert.Equal(t, 14, account.Balance)
	}
	os.Remove("flux_smoke.db")

}
