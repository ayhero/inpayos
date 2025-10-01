package services

import (
	"inpayos/internal/protocol"
	"sync"
)

type TransactionService struct {
	PayinService  *PayinService
	PayoutService *PayoutService
}

var (
	transactionService     *TransactionService
	transactionServiceOnce sync.Once
)

func SetupTransactionService() {
	transactionServiceOnce.Do(func() {
		transactionService = &TransactionService{
			PayinService:  GetPayinService(),
			PayoutService: GetPayoutService(),
		}
	})
}

// GetTransactionService 获取Transaction服务单例
func GetTransactionService() *TransactionService {
	if transactionService == nil {
		SetupTransactionService()
	}
	return transactionService
}

func (s *TransactionService) CreatePayin(req *protocol.CreateTransactionRequest) (trx *protocol.Transaction, code protocol.ErrorCode) {
	return &protocol.Transaction{}, protocol.Success
}

func (s *TransactionService) CreatePayout(req *protocol.CreateTransactionRequest) (trx *protocol.Transaction, code protocol.ErrorCode) {
	return &protocol.Transaction{}, protocol.Success
}

func (s *TransactionService) Cancel(req *protocol.CancelTransactionRequest) (trx *protocol.Transaction, code protocol.ErrorCode) {
	return &protocol.Transaction{}, protocol.Success
}
func (s *TransactionService) Query(req *protocol.QueryTransactionRequest) (trx *protocol.Transaction, code protocol.ErrorCode) {
	return &protocol.Transaction{}, protocol.Success
}
