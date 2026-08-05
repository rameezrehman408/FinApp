package investment

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AssetType string

const (
	AssetTypeStock      AssetType = "STOCK"
	AssetTypeBond       AssetType = "BOND"
	AssetTypeMutualFund AssetType = "MUTUAL_FUND"
	AssetTypeETF        AssetType = "ETF"
	AssetTypeCrypto     AssetType = "CRYPTO"
	AssetTypeRealEstate AssetType = "REAL_ESTATE"
	AssetTypeCommodity  AssetType = "COMMODITY"
	AssetTypeOther      AssetType = "OTHER"
)

type Asset struct {
	ID          uuid.UUID
	Symbol      string
	Name        string
	Type        AssetType
	Currency    string
	Exchange    string
	ISIN        string
	CUSIP       string
	Sector      string
	Industry    string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int
}

func NewAsset(symbol, name string, assetType AssetType, currency string) *Asset {
	now := time.Now().UTC()
	return &Asset{
		ID:       uuid.New(),
		Symbol:   symbol,
		Name:     name,
		Type:     assetType,
		Currency: currency,
		IsActive: true,
		CreatedAt: now,
		UpdatedAt: now,
		Version:   1,
	}
}

type TransactionType string

const (
	TransactionTypeBuy    TransactionType = "BUY"
	TransactionTypeSell   TransactionType = "SELL"
	TransactionTypeDividend TransactionType = "DIVIDEND"
	TransactionTypeSplit  TransactionType = "SPLIT"
	TransactionTypeMerge  TransactionType = "MERGE"
	TransactionTypeSpinOff TransactionType = "SPIN_OFF"
)

type InvestmentTransaction struct {
	ID              uuid.UUID
	AccountID       uuid.UUID
	AssetID         uuid.UUID
	Type            TransactionType
	Quantity        decimal.Decimal
	PricePerUnit    decimal.Decimal
	Fees            decimal.Decimal
	Taxes           decimal.Decimal
	TotalAmount     decimal.Decimal
	Currency        string
	TransactionDate time.Time
	SettlementDate  time.Time
	Notes           string
	Reference       string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Version         int
}

func NewBuyTransaction(accountID, assetID uuid.UUID, quantity, pricePerUnit, fees, taxes decimal.Decimal, currency string, transactionDate, settlementDate time.Time) *InvestmentTransaction {
	total := quantity.Mul(pricePerUnit).Add(fees).Add(taxes)
	now := time.Now().UTC()
	return &InvestmentTransaction{
		ID:              uuid.New(),
		AccountID:       accountID,
		AssetID:         assetID,
		Type:            TransactionTypeBuy,
		Quantity:        quantity,
		PricePerUnit:    pricePerUnit,
		Fees:            fees,
		Taxes:           taxes,
		TotalAmount:     total,
		Currency:        currency,
		TransactionDate: transactionDate,
		SettlementDate:  settlementDate,
		CreatedAt:       now,
		UpdatedAt:       now,
		Version:         1,
	}
}

func NewSellTransaction(accountID, assetID uuid.UUID, quantity, pricePerUnit, fees, taxes decimal.Decimal, currency string, transactionDate, settlementDate time.Time) *InvestmentTransaction {
	total := quantity.Mul(pricePerUnit).Sub(fees).Sub(taxes)
	now := time.Now().UTC()
	return &InvestmentTransaction{
		ID:              uuid.New(),
		AccountID:       accountID,
		AssetID:         assetID,
		Type:            TransactionTypeSell,
		Quantity:        quantity,
		PricePerUnit:    pricePerUnit,
		Fees:            fees,
		Taxes:           taxes,
		TotalAmount:     total,
		Currency:        currency,
		TransactionDate: transactionDate,
		SettlementDate:  settlementDate,
		CreatedAt:       now,
		UpdatedAt:       now,
		Version:         1,
	}
}

type Holding struct {
	ID              uuid.UUID
	AccountID       uuid.UUID
	AssetID         uuid.UUID
	Quantity        decimal.Decimal
	AverageCost     decimal.Decimal
	TotalCost       decimal.Decimal
	CurrentPrice    decimal.Decimal
	CurrentValue    decimal.Decimal
	UnrealizedGain  decimal.Decimal
	RealizedGain    decimal.Decimal
	Currency        string
	FirstPurchaseDate time.Time
	LastTransactionDate time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Version         int
}

func NewHolding(accountID, assetID uuid.UUID, quantity, averageCost, totalCost decimal.Decimal, currency string, firstPurchaseDate time.Time) *Holding {
	now := time.Now().UTC()
	return &Holding{
		ID:                 uuid.New(),
		AccountID:          accountID,
		AssetID:            assetID,
		Quantity:           quantity,
		AverageCost:        averageCost,
		TotalCost:          totalCost,
		Currency:           currency,
		FirstPurchaseDate:  firstPurchaseDate,
		LastTransactionDate: firstPurchaseDate,
		CreatedAt:          now,
		UpdatedAt:          now,
		Version:            1,
	}
}

func (h *Holding) AddPurchase(quantity, pricePerUnit, fees, taxes decimal.Decimal, transactionDate time.Time) {
	totalCost := quantity.Mul(pricePerUnit).Add(fees).Add(taxes)
	newQuantity := h.Quantity.Add(quantity)
	newTotalCost := h.TotalCost.Add(totalCost)
	h.AverageCost = newTotalCost.Div(newQuantity)
	h.Quantity = newQuantity
	h.TotalCost = newTotalCost
	h.LastTransactionDate = transactionDate
	h.UpdatedAt = time.Now().UTC()
	h.Version++
}

func (h *Holding) RemoveSale(quantity, pricePerUnit, fees, taxes decimal.Decimal, transactionDate time.Time) decimal.Decimal {
	costBasis := quantity.Mul(h.AverageCost)
	proceeds := quantity.Mul(pricePerUnit).Sub(fees).Sub(taxes)
	realizedGain := proceeds.Sub(costBasis)

	h.RealizedGain = h.RealizedGain.Add(realizedGain)
	h.Quantity = h.Quantity.Sub(quantity)
	h.TotalCost = h.TotalCost.Sub(costBasis)
	h.LastTransactionDate = transactionDate
	h.UpdatedAt = time.Now().UTC()
	h.Version++

	return realizedGain
}

func (h *Holding) UpdateCurrentPrice(price decimal.Decimal) {
	h.CurrentPrice = price
	h.CurrentValue = h.Quantity.Mul(price)
	h.UnrealizedGain = h.CurrentValue.Sub(h.TotalCost)
	h.UpdatedAt = time.Now().UTC()
	h.Version++
}

func (h *Holding) TotalGain() decimal.Decimal {
	return h.UnrealizedGain.Add(h.RealizedGain)
}

type CapitalGain struct {
	ID              uuid.UUID
	AccountID       uuid.UUID
	AssetID         uuid.UUID
	TransactionID   uuid.UUID
	Type            CapitalGainType
	Quantity        decimal.Decimal
	CostBasis       decimal.Decimal
	Proceeds        decimal.Decimal
	GainLoss        decimal.Decimal
	HoldingPeriod   int
	IsShortTerm     bool
	TaxYear         int
	TransactionDate time.Time
	CreatedAt       time.Time
}

type CapitalGainType string

const (
	CapitalGainTypeShortTerm CapitalGainType = "SHORT_TERM"
	CapitalGainTypeLongTerm  CapitalGainType = "LONG_TERM"
)

func NewCapitalGain(accountID, assetID, transactionID uuid.UUID, quantity, costBasis, proceeds decimal.Decimal, transactionDate time.Time, purchaseDate time.Time, taxYear int) *CapitalGain {
	holdingPeriod := int(transactionDate.Sub(purchaseDate).Hours() / 24)
	isShortTerm := holdingPeriod <= 365
	gainLoss := proceeds.Sub(costBasis)
	gainType := CapitalGainTypeLongTerm
	if isShortTerm {
		gainType = CapitalGainTypeShortTerm
	}

	now := time.Now().UTC()
	return &CapitalGain{
		ID:              uuid.New(),
		AccountID:       accountID,
		AssetID:         assetID,
		TransactionID:   transactionID,
		Type:            gainType,
		Quantity:        quantity,
		CostBasis:       costBasis,
		Proceeds:        proceeds,
		GainLoss:        gainLoss,
		HoldingPeriod:   holdingPeriod,
		IsShortTerm:     isShortTerm,
		TaxYear:         taxYear,
		TransactionDate: transactionDate,
		CreatedAt:       now,
	}
}