package protocol

import "github.com/shopspring/decimal"

type StatisModel struct {
	Price *decimal.Decimal `json:"price" gorm:"column:price;type:numeric"`
	Total int              `json:"total" gorm:"column:total;type:int"`
}
