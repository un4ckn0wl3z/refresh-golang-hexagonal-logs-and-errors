package service

import (
	"bank/errs"
	"bank/logs"
	"bank/repository"
	"strings"
	"time"
)

type accountServiceImpl struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return accountServiceImpl{accountRepo: accountRepo}
}

func (s accountServiceImpl) NewAccount(cusID int, accReq NewAccountRequest) (*AccountResponse, error) {

	if accReq.Amount < 5000 {
		return nil, errs.NewValidationError("amount at least 5,000")
	}

	if strings.ToLower(accReq.AccountType) != "checking" && strings.ToLower(accReq.AccountType) != "saving" {
		return nil, errs.NewValidationError("account type should be saving or checking")
	}

	account := repository.Account{
		CustomerID:  cusID,
		OpeningDate: time.Now().Format("2006-1-2 15:04:05"),
		AccountType: accReq.AccountType,
		Amount:      accReq.Amount,
		Status:      1,
	}
	newAcc, err := s.accountRepo.Create(account)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	response := AccountResponse{
		AccountID:   newAcc.AccountID,
		OpeningDate: newAcc.OpeningDate,
		AccountType: newAcc.AccountType,
		Amount:      newAcc.Amount,
		Status:      newAcc.Status,
	}

	return &response, nil

}

func (s accountServiceImpl) GetAccounts(cusID int) ([]AccountResponse, error) {
	accounts, err := s.accountRepo.GetAll(cusID)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	responses := []AccountResponse{}
	for _, account := range accounts {
		responses = append(responses, AccountResponse{
			AccountID:   account.AccountID,
			OpeningDate: account.OpeningDate,
			AccountType: account.AccountType,
			Amount:      account.Amount,
			Status:      account.Status,
		})
	}

	return responses, nil

}
