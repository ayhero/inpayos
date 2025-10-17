package protocol

// 扩展货币代码常量 (ISO 4217) - 补充 account.go 中未定义的货币
const (
	CcyINR = "INR" // 印度卢比
	CcyKRW = "KRW" // 韩元
	CcyAUD = "AUD" // 澳元
	CcyCAD = "CAD" // 加元
	CcySGD = "SGD" // 新加坡元
	CcyHKD = "HKD" // 港币
	CcyTHB = "THB" // 泰铢
	CcyMYR = "MYR" // 马来西亚林吉特
	CcyIDR = "IDR" // 印尼盾
	CcyPHP = "PHP" // 菲律宾比索
	CcyVND = "VND" // 越南盾
)

// 加密货币常量
const (
	CcyUSDT = "USDT" // 泰达币
	CcyBTC  = "BTC"  // 比特币
	CcyETH  = "ETH"  // 以太坊
)

// 货币信息结构
type CurrencyInfo struct {
	Code         string   `json:"code"`      // 货币代码
	Name         string   `json:"name"`      // 货币名称
	Symbol       string   `json:"symbol"`    // 货币符号
	Decimals     int      `json:"decimals"`  // 小数位数
	CountryCodes []string `json:"countries"` // 主要使用国家代码
}

// 货币信息映射
var CurrencyInfoMap = map[string]CurrencyInfo{
	CcyUSD: {
		Code:         CcyUSD,
		Name:         "US Dollar",
		Symbol:       "$",
		Decimals:     2,
		CountryCodes: []string{"USA"},
	},
	CcyEUR: {
		Code:         CcyEUR,
		Name:         "Euro",
		Symbol:       "€",
		Decimals:     2,
		CountryCodes: []string{"DEU", "FRA", "ITA", "ESP"},
	},
	CcyGBP: {
		Code:         CcyGBP,
		Name:         "British Pound",
		Symbol:       "£",
		Decimals:     2,
		CountryCodes: []string{"GBR"},
	},
	CcyJPY: {
		Code:         CcyJPY,
		Name:         "Japanese Yen",
		Symbol:       "¥",
		Decimals:     0,
		CountryCodes: []string{"JPN"},
	},
	CcyCNY: {
		Code:         CcyCNY,
		Name:         "Chinese Yuan",
		Symbol:       "¥",
		Decimals:     2,
		CountryCodes: []string{"CHN"},
	},
	CcyINR: {
		Code:         CcyINR,
		Name:         "Indian Rupee",
		Symbol:       "₹",
		Decimals:     2,
		CountryCodes: []string{"IND"},
	},
	CcyKRW: {
		Code:         CcyKRW,
		Name:         "South Korean Won",
		Symbol:       "₩",
		Decimals:     0,
		CountryCodes: []string{"KOR"},
	},
	CcyAUD: {
		Code:         CcyAUD,
		Name:         "Australian Dollar",
		Symbol:       "A$",
		Decimals:     2,
		CountryCodes: []string{"AUS"},
	},
	CcyCAD: {
		Code:         CcyCAD,
		Name:         "Canadian Dollar",
		Symbol:       "C$",
		Decimals:     2,
		CountryCodes: []string{"CAN"},
	},
	CcySGD: {
		Code:         CcySGD,
		Name:         "Singapore Dollar",
		Symbol:       "S$",
		Decimals:     2,
		CountryCodes: []string{"SGP"},
	},
	CcyHKD: {
		Code:         CcyHKD,
		Name:         "Hong Kong Dollar",
		Symbol:       "HK$",
		Decimals:     2,
		CountryCodes: []string{"HKG"},
	},
	CcyTHB: {
		Code:         CcyTHB,
		Name:         "Thai Baht",
		Symbol:       "฿",
		Decimals:     2,
		CountryCodes: []string{"THA"},
	},
	CcyMYR: {
		Code:         CcyMYR,
		Name:         "Malaysian Ringgit",
		Symbol:       "RM",
		Decimals:     2,
		CountryCodes: []string{"MYS"},
	},
	CcyIDR: {
		Code:         CcyIDR,
		Name:         "Indonesian Rupiah",
		Symbol:       "Rp",
		Decimals:     0,
		CountryCodes: []string{"IDN"},
	},
	CcyPHP: {
		Code:         CcyPHP,
		Name:         "Philippine Peso",
		Symbol:       "₱",
		Decimals:     2,
		CountryCodes: []string{"PHL"},
	},
	CcyVND: {
		Code:         CcyVND,
		Name:         "Vietnamese Dong",
		Symbol:       "₫",
		Decimals:     0,
		CountryCodes: []string{"VNM"},
	},
	CcyUSDT: {
		Code:         CcyUSDT,
		Name:         "Tether USD",
		Symbol:       "USDT",
		Decimals:     6,
		CountryCodes: []string{},
	},
	CcyBTC: {
		Code:         CcyBTC,
		Name:         "Bitcoin",
		Symbol:       "₿",
		Decimals:     8,
		CountryCodes: []string{},
	},
	CcyETH: {
		Code:         CcyETH,
		Name:         "Ethereum",
		Symbol:       "Ξ",
		Decimals:     18,
		CountryCodes: []string{},
	},
}

// 货币验证函数
func IsValidCurrency(ccy string) bool {
	_, exists := CurrencyInfoMap[ccy]
	return exists
}

// 获取货币信息
func GetCurrencyInfo(ccy string) (*CurrencyInfo, bool) {
	info, exists := CurrencyInfoMap[ccy]
	if exists {
		return &info, true
	}
	return nil, false
}

// 获取所有支持的货币代码
func GetSupportedCurrencies() []string {
	currencies := make([]string, 0, len(CurrencyInfoMap))
	for code := range CurrencyInfoMap {
		currencies = append(currencies, code)
	}
	return currencies
}

// 按区域获取货币
func GetCurrenciesByRegion() map[string][]string {
	return map[string][]string{
		"Asia":     {CcyINR, CcyJPY, CcyCNY, CcyKRW, CcySGD, CcyHKD, CcyTHB, CcyMYR, CcyIDR, CcyPHP, CcyVND},
		"Europe":   {CcyEUR, CcyGBP},
		"Americas": {CcyUSD, CcyCAD},
		"Oceania":  {CcyAUD},
		"Crypto":   {CcyUSDT, CcyBTC, CcyETH},
	}
}
