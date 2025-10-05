package services

import (
	"inpayos/internal/protocol"
	"sync"
)

type MerchantTransactionService struct {
	PayinService  *MerchantPayinService
	PayoutService *MerchantPayoutService
}

var (
	transactionService     *MerchantTransactionService
	transactionServiceOnce sync.Once
)

func SetupTransactionService() {
	transactionServiceOnce.Do(func() {
		transactionService = &MerchantTransactionService{
			PayinService:  GetMerchantPayinService(),
			PayoutService: GetMerchantPayoutService(),
		}
	})
}

// GetMerchantTransactionService 获取Transaction服务单例
func GetMerchantTransactionService() *MerchantTransactionService {
	if transactionService == nil {
		SetupTransactionService()
	}
	return transactionService
}

func (s *MerchantTransactionService) CreatePayin(req *protocol.CreateTransactionRequest) (trx *protocol.Transaction, code protocol.ErrorCode) {
	return &protocol.Transaction{}, protocol.Success
}

func (s *MerchantTransactionService) CreatePayout(req *protocol.CreateTransactionRequest) (trx *protocol.Transaction, code protocol.ErrorCode) {
	return &protocol.Transaction{}, protocol.Success
}

func (s *MerchantTransactionService) Cancel(req *protocol.CancelTransactionRequest) (trx *protocol.Transaction, code protocol.ErrorCode) {
	return &protocol.Transaction{}, protocol.Success
}
func (s *MerchantTransactionService) Query(req *protocol.QueryTransactionRequest) (trx *protocol.Transaction, code protocol.ErrorCode) {
	return &protocol.Transaction{}, protocol.Success
}
