package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/seisram/FinApp/internal/application/command"
	"github.com/seisram/FinApp/internal/application/query"
)

type Handlers struct {
	CreateAccount      *command.CreateAccountHandler
	CreateTransaction  *command.CreateTransactionHandler
	CreateCategory     *command.CreateCategoryHandler
	CreateParty        *command.CreatePartyHandler
	AccountQuery       *query.AccountQueryHandler
	TransactionQuery   *query.TransactionQueryHandler
	CategoryQuery      *query.CategoryQueryHandler
	PartyQuery         *query.PartyQueryHandler
}

func SetupRoutes(r *gin.Engine, h Handlers) {
	api := r.Group("/api/v1")
	{
		// Health check
		api.GET("/health", healthCheck)

		// Accounts
		accounts := api.Group("/accounts")
		{
			accounts.POST("", h.createAccount)
			accounts.GET("", h.listAccounts)
			accounts.GET("/:id", h.getAccount)
		}

		// Transactions
		transactions := api.Group("/transactions")
		{
			transactions.POST("", h.createTransaction)
			transactions.GET("", h.listTransactions)
			transactions.GET("/:id", h.getTransaction)
		}

		// Categories
		categories := api.Group("/categories")
		{
			categories.POST("", h.createCategory)
			categories.GET("", h.listCategories)
			categories.GET("/:id", h.getCategory)
		}

		// Parties
		parties := api.Group("/parties")
		{
			parties.POST("", h.createParty)
			parties.GET("", h.listParties)
			parties.GET("/:id", h.getParty)
		}
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "finapp-api",
	})
}

// Account handlers
func (h *Handlers) createAccount(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Currency    string `json:"currency" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := command.CreateAccountCommand{
		Name:        req.Name,
		Type:        command.AccountType(req.Type),
		Currency:    req.Currency,
		Description: req.Description,
	}

	result, err := h.CreateAccount.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *Handlers) getAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	result, err := h.AccountQuery.HandleGetAccount(c.Request.Context(), query.GetAccountQuery{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result == nil || result.Account == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handlers) listAccounts(c *gin.Context) {
	accountType := c.Query("type")
	var accType *command.AccountType
	if accountType != "" {
		t := command.AccountType(accountType)
		accType = &t
	}

	result, err := h.AccountQuery.HandleListAccounts(c.Request.Context(), query.ListAccountsQuery{Type: accType})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Transaction handlers
func (h *Handlers) createTransaction(c *gin.Context) {
	var req struct {
		AccountID       string  `json:"account_id" binding:"required"`
		CounterAccountID *string `json:"counter_account_id"`
		CategoryID      *string `json:"category_id"`
		PartyID         *string `json:"party_id"`
		Type            string  `json:"type" binding:"required"`
		Amount          string  `json:"amount" binding:"required"`
		Currency        string  `json:"currency" binding:"required"`
		Description     string  `json:"description"`
		Reference       string  `json:"reference"`
		TransactionDate string  `json:"transaction_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
		return
	}

	var counterAccountID *uuid.UUID
	if req.CounterAccountID != nil {
		id, err := uuid.Parse(*req.CounterAccountID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid counter_account_id"})
			return
		}
		counterAccountID = &id
	}

	var categoryID *uuid.UUID
	if req.CategoryID != nil {
		id, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id"})
			return
		}
		categoryID = &id
	}

	var partyID *uuid.UUID
	if req.PartyID != nil {
		id, err := uuid.Parse(*req.PartyID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid party_id"})
			return
		}
		partyID = &id
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid amount"})
		return
	}

	transactionDate, err := time.Parse(time.RFC3339, req.TransactionDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction_date format (use RFC3339)"})
		return
	}

	cmd := command.CreateTransactionCommand{
		AccountID:       accountID,
		CounterAccountID: counterAccountID,
		CategoryID:      categoryID,
		PartyID:         partyID,
		Type:            command.TransactionType(req.Type),
		Amount:          amount,
		Currency:        req.Currency,
		Description:     req.Description,
		Reference:       req.Reference,
		TransactionDate: transactionDate,
	}

	result, err := h.CreateTransaction.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *Handlers) getTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	result, err := h.TransactionQuery.HandleGetTransaction(c.Request.Context(), query.GetTransactionQuery{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result == nil || result.Transaction == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handlers) listTransactions(c *gin.Context) {
	accountIDStr := c.Query("account_id")
	categoryIDStr := c.Query("category_id")
	partyIDStr := c.Query("party_id")
	statusStr := c.Query("status")
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	var accountID *uuid.UUID
	if accountIDStr != "" {
		id, err := uuid.Parse(accountIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
			return
		}
		accountID = &id
	}

	var categoryID *uuid.UUID
	if categoryIDStr != "" {
		id, err := uuid.Parse(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category_id"})
			return
		}
		categoryID = &id
	}

	var partyID *uuid.UUID
	if partyIDStr != "" {
		id, err := uuid.Parse(partyIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid party_id"})
			return
		}
		partyID = &id
	}

	var status *command.TransactionStatus
	if statusStr != "" {
		s := command.TransactionStatus(statusStr)
		status = &s
	}

	result, err := h.TransactionQuery.HandleListTransactions(c.Request.Context(), query.ListTransactionsQuery{
		AccountID:  accountID,
		CategoryID: categoryID,
		PartyID:    partyID,
		Status:     status,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Category handlers
func (h *Handlers) createCategory(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Type        string  `json:"type" binding:"required"`
		ParentID    *string `json:"parent_id"`
		Description string  `json:"description"`
		Icon        string  `json:"icon"`
		Color       string  `json:"color"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var parentID *uuid.UUID
	if req.ParentID != nil {
		id, err := uuid.Parse(*req.ParentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent_id"})
			return
		}
		parentID = &id
	}

	cmd := command.CreateCategoryCommand{
		Name:        req.Name,
		Type:        command.CategoryType(req.Type),
		ParentID:    parentID,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
	}

	result, err := h.CreateCategory.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *Handlers) getCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	result, err := h.CategoryQuery.HandleGetCategory(c.Request.Context(), query.GetCategoryQuery{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result == nil || result.Category == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handlers) listCategories(c *gin.Context) {
	typeStr := c.Query("type")
	parentIDStr := c.Query("parent_id")
	rootOnlyStr := c.Query("root_only")

	var catType *command.CategoryType
	if typeStr != "" {
		t := command.CategoryType(typeStr)
		catType = &t
	}

	var parentID *uuid.UUID
	if parentIDStr != "" {
		id, err := uuid.Parse(parentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent_id"})
			return
		}
		parentID = &id
	}

	rootOnly := rootOnlyStr == "true"

	result, err := h.CategoryQuery.HandleListCategories(c.Request.Context(), query.ListCategoriesQuery{
		Type:     catType,
		ParentID: parentID,
		RootOnly: rootOnly,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Party handlers
func (h *Handlers) createParty(c *gin.Context) {
	var req struct {
		Name          string  `json:"name" binding:"required"`
		Type          string  `json:"type" binding:"required"`
		TaxID         string  `json:"tax_id"`
		Email         string  `json:"email"`
		Phone         string  `json:"phone"`
		Street        string  `json:"street"`
		City          string  `json:"city"`
		State         string  `json:"state"`
		PostalCode    string  `json:"postal_code"`
		Country       string  `json:"country"`
		ContactPerson string  `json:"contact_person"`
		Notes         string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := command.CreatePartyCommand{
		Name:          req.Name,
		Type:          command.PartyType(req.Type),
		TaxID:         req.TaxID,
		Email:         req.Email,
		Phone:         req.Phone,
		Street:        req.Street,
		City:          req.City,
		State:         req.State,
		PostalCode:    req.PostalCode,
		Country:       req.Country,
		ContactPerson: req.ContactPerson,
		Notes:         req.Notes,
	}

	result, err := h.CreateParty.Handle(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *Handlers) getParty(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid party ID"})
		return
	}

	result, err := h.PartyQuery.HandleGetParty(c.Request.Context(), query.GetPartyQuery{ID: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result == nil || result.Party == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "party not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handlers) listParties(c *gin.Context) {
	typeStr := c.Query("type")
	activeStr := c.Query("active")
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	var partyType *command.PartyType
	if typeStr != "" {
		t := command.PartyType(typeStr)
		partyType = &t
	}

	var active *bool
	if activeStr != "" {
		a := activeStr == "true"
		active = &a
	}

	result, err := h.PartyQuery.HandleListParties(c.Request.Context(), query.ListPartiesQuery{
		Type:   partyType,
		Active: active,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}