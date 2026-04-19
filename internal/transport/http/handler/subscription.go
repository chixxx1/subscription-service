package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chixxx1/subscription-service/internal/domain"
	sub_service "github.com/chixxx1/subscription-service/internal/service/subscription"
	"github.com/chixxx1/subscription-service/internal/transport/http/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	service *sub_service.SubscriptionService
	logger  *zap.Logger
}

func NewSubscriptionHandler(svc *sub_service.SubscriptionService, logger *zap.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: svc,
		logger:  logger,
	}
}

// CreateSubscription godoc
// @Summary      Create a new subscription
// @Description  Create a new subscription record for a user
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body      dto.CreateSubscriptionRequest  true  "Subscription data"
// @Success      201           {object}  map[string]string
// @Failure      400           {object}  dto.ErrorResponse
// @Failure      500           {object}  dto.ErrorResponse
// @Router       /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	_, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id format, must be valid UUID"})
		return
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid start_date format, use MM-YYYY"})
		return
	}

	var endDate *time.Time
	if req.EndDate != "" {
		end, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid end_date format, use MM-YYYY"})
			return
		}
		endDate = &end
	}

	sub := domain.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := h.service.Create(c.Request.Context(), sub); err != nil {
		h.logger.Error("failed to create subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

// GetByID godoc
// @Summary			 Get subscription by ID
// @Description	 Retrieve details of a special subscription by its ID
// @Tags         subscriptions
// @Produce      json
// @Param        id  path      int  true  "Subscription ID"
// @Success      200  {object}  domain.Subscription
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id format"})
		return
	}

	sub, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// List godoc
// @Summary      List all subscriptions
// @Description  Retrieve a list of subscriptions with optional filtering and pagination
// @Tags         subscriptions
// @Produce      json
// @Param        user_id       query     string  false  "Filter by User ID (UUID)"
// @Param        service_name  query     string  false  "Filter by Service Name"
// @Param        limit         query     int     false  "Number of records to return (default 50)"
// @Param        offset        query     int     false  "Number of records to skip (default 0)"
// @Success      200           {array}   domain.Subscription
// @Failure      400           {object}  dto.ErrorResponse
// @Failure      500           {object}  dto.ErrorResponse
// @Router       /subscriptions [get]
func (h *SubscriptionHandler) List(c *gin.Context) {
	var query dto.ListSubscriptionsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	filter := domain.SubscriptionFilter{
		UserID:      query.UserID,
		ServiceName: query.ServiceName,
		Limit:       query.Limit,
		Offset:      query.Offset,
	}

	subs, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("failed to list subscriptions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, subs)
}

// Update godoc
// @Summary      Update an existing subscription
// @Description  Update details of a subscription by its ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id            path      int                          true  "Subscription ID"
// @Param        subscription  body      dto.UpdateSubscriptionRequest  true  "Updated subscription data"
// @Success      200           {object}  map[string]string
// @Failure      400           {object}  dto.ErrorResponse
// @Failure      500           {object}  dto.ErrorResponse
// @Router       /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id format"})
		return
	}

	var req dto.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	_, err = uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id format, must be valid UUID"})
		return
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid start_date format, use MM-YYYY"})
		return
	}

	var endDate *time.Time
	if req.EndDate != "" {
		end, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid end_date format, use MM-YYYY"})
			return
		}
		endDate = &end
	}

	sub := domain.Subscription{
		ID:          id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := h.service.Update(c.Request.Context(), sub); err != nil {
		h.logger.Error("failed to update subscription", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// Delete godoc
// @Summary      Delete a subscription
// @Description  Remove a subscription record by its ID
// @Tags         subscriptions
// @Produce      json
// @Param        id  path      int  true  "Subscription ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id format"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete subscription", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// GetTotalCost godoc
// @Summary      Get total cost of subscriptions
// @Description  Calculate total cost for a period with optional filters
// @Tags         subscriptions
// @Produce      json
// @Param        user_id       query     string  false  "User ID (UUID)"
// @Param        service_name  query     string  false  "Service Name"
// @Param        period_start  query     string  true   "Start Period (YYYY-MM)"
// @Param        period_end    query     string  true   "End Period (YYYY-MM)"
// @Success      200           {object}  map[string]int64
// @Failure      400           {object}  dto.ErrorResponse
// @Failure      500           {object}  dto.ErrorResponse
// @Router       /subscriptions/total-cost [get]
func (h *SubscriptionHandler) GetTotalCost(c *gin.Context) {
	var req dto.TotalCostQuery

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01", req.PeriodStart)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid period_start format, use YYYY-MM"})
		return
	}

	endDate, err := time.Parse("2006-01", req.PeriodEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid period_end format, use YYYY-MM"})
		return
	}

	lastDayOfMonth := endDate.AddDate(0, 1, -1)

	domainReq := domain.TotalCostRequest{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		PeriodStart: startDate,
		PeriodEnd:   lastDayOfMonth,
	}

	total, err := h.service.GetTotalCost(c.Request.Context(), domainReq)
	if err != nil {
		h.logger.Error("failed to get total cost", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
	}

	c.JSON(http.StatusOK, gin.H{"total_price": total})
}
