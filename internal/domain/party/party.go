package party

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PartyType string

const (
	PartyTypePerson    PartyType = "PERSON"
	PartyTypeCompany   PartyType = "COMPANY"
	PartyTypeBank      PartyType = "BANK"
	PartyTypeGovernment PartyType = "GOVERNMENT"
	PartyTypeOther     PartyType = "OTHER"
)

type Party struct {
	ID           uuid.UUID
	Name         string
	Type         PartyType
	TaxID        string
	Email        string
	Phone        string
	Address      Address
	ContactPerson string
	Notes        string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Version      int
}

type Address struct {
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}

func NewParty(name string, partyType PartyType) *Party {
	now := time.Now().UTC()
	return &Party{
		ID:        uuid.New(),
		Name:      name,
		Type:      partyType,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
		Version:   1,
	}
}

func (p *Party) Update(name string, email string, phone string, address Address) {
	p.Name = name
	p.Email = email
	p.Phone = phone
	p.Address = address
	p.UpdatedAt = time.Now().UTC()
	p.Version++
}

func (p *Party) Deactivate() {
	p.IsActive = false
	p.UpdatedAt = time.Now().UTC()
	p.Version++
}

type Receivable struct {
	ID              uuid.UUID
	PartyID         uuid.UUID
	Description     string
	Amount          decimal.Decimal
	Currency        string
	DueDate         time.Time
	Status          ReceivableStatus
	InvoiceNumber   string
	InvoiceDate     time.Time
	PaidAmount      decimal.Decimal
	PaidDate        *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Version         int
}

type ReceivableStatus string

const (
	ReceivableStatusPending   ReceivableStatus = "PENDING"
	ReceivableStatusPartial   ReceivableStatus = "PARTIAL"
	ReceivableStatusPaid      ReceivableStatus = "PAID"
	ReceivableStatusOverdue   ReceivableStatus = "OVERDUE"
	ReceivableStatusCancelled ReceivableStatus = "CANCELLED"
)

func NewReceivable(partyID uuid.UUID, description string, amount decimal.Decimal, currency string, dueDate time.Time, invoiceNumber string) *Receivable {
	now := time.Now().UTC()
	return &Receivable{
		ID:            uuid.New(),
		PartyID:       partyID,
		Description:   description,
		Amount:        amount,
		Currency:      currency,
		DueDate:       dueDate,
		Status:        ReceivableStatusPending,
		InvoiceNumber: invoiceNumber,
		InvoiceDate:   now,
		PaidAmount:    decimal.Zero,
		CreatedAt:     now,
		UpdatedAt:     now,
		Version:       1,
	}
}

func (r *Receivable) RecordPayment(amount decimal.Decimal, paidDate time.Time) error {
	if r.Status == ReceivableStatusPaid {
		return ErrAlreadyPaid
	}
	if r.Status == ReceivableStatusCancelled {
		return ErrCancelled
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}
	if r.PaidAmount.Add(amount).GreaterThan(r.Amount) {
		return ErrOverpayment
	}

	r.PaidAmount = r.PaidAmount.Add(amount)
	r.PaidDate = &paidDate
	r.UpdatedAt = time.Now().UTC()
	r.Version++

	if r.PaidAmount.Equal(r.Amount) {
		r.Status = ReceivableStatusPaid
	} else {
		r.Status = ReceivableStatusPartial
	}
	return nil
}

func (r *Receivable) OutstandingAmount() decimal.Decimal {
	return r.Amount.Sub(r.PaidAmount)
}

func (r *Receivable) IsOverdue() bool {
	if r.Status == ReceivableStatusPaid || r.Status == ReceivableStatusCancelled {
		return false
	}
	return time.Now().After(r.DueDate)
}

type Payable struct {
	ID              uuid.UUID
	PartyID         uuid.UUID
	Description     string
	Amount          decimal.Decimal
	Currency        string
	DueDate         time.Time
	Status          PayableStatus
	InvoiceNumber   string
	InvoiceDate     time.Time
	PaidAmount      decimal.Decimal
	PaidDate        *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Version         int
}

type PayableStatus string

const (
	PayableStatusPending   PayableStatus = "PENDING"
	PayableStatusPartial   PayableStatus = "PARTIAL"
	PayableStatusPaid      PayableStatus = "PAID"
	PayableStatusOverdue   PayableStatus = "OVERDUE"
	PayableStatusCancelled PayableStatus = "CANCELLED"
)

func NewPayable(partyID uuid.UUID, description string, amount decimal.Decimal, currency string, dueDate time.Time, invoiceNumber string) *Payable {
	now := time.Now().UTC()
	return &Payable{
		ID:            uuid.New(),
		PartyID:       partyID,
		Description:   description,
		Amount:        amount,
		Currency:      currency,
		DueDate:       dueDate,
		Status:        PayableStatusPending,
		InvoiceNumber: invoiceNumber,
		InvoiceDate:   now,
		PaidAmount:    decimal.Zero,
		CreatedAt:     now,
		UpdatedAt:     now,
		Version:       1,
	}
}

func (p *Payable) RecordPayment(amount decimal.Decimal, paidDate time.Time) error {
	if p.Status == PayableStatusPaid {
		return ErrAlreadyPaid
	}
	if p.Status == PayableStatusCancelled {
		return ErrCancelled
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}
	if p.PaidAmount.Add(amount).GreaterThan(p.Amount) {
		return ErrOverpayment
	}

	p.PaidAmount = p.PaidAmount.Add(amount)
	p.PaidDate = &paidDate
	p.UpdatedAt = time.Now().UTC()
	p.Version++

	if p.PaidAmount.Equal(p.Amount) {
		p.Status = PayableStatusPaid
	} else {
		p.Status = PayableStatusPartial
	}
	return nil
}

func (p *Payable) OutstandingAmount() decimal.Decimal {
	return p.Amount.Sub(p.PaidAmount)
}

func (p *Payable) IsOverdue() bool {
	if p.Status == PayableStatusPaid || p.Status == PayableStatusCancelled {
		return false
	}
	return time.Now().After(p.DueDate)
}

var (
	ErrAlreadyPaid   = &DomainError{Code: "ALREADY_PAID", Message: "Already fully paid"}
	ErrCancelled     = &DomainError{Code: "CANCELLED", Message: "Document is cancelled"}
	ErrInvalidAmount = &DomainError{Code: "INVALID_AMOUNT", Message: "Invalid payment amount"}
	ErrOverpayment   = &DomainError{Code: "OVERPAYMENT", Message: "Payment exceeds outstanding amount"}
)

type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}