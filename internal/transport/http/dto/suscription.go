package dto

import "github.com/gin-gonic/gin"

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required,min=1"`
	UserID      string `json:"user_id" binding:"required,uuid"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required,min=1"`
	UserID      string `json:"user_id" binding:"required,uuid"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date,omitempty"`
}

type ListSubscriptionsQuery struct {
	UserID      string `form:"user_id"`
	ServiceName string `form:"service_name"`
	Limit       int    `form:"limit" default:"50"`
	Offset      int    `form:"offset" default:"0"`
}

type TotalCostQuery struct {
	UserID      string `form:"user_id"`
	ServiceName string `form:"service_name"`
	PeriodStart string `form:"period_start" binding:"required"`
	PeriodEnd   string `form:"period_end" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

func bindError(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(400, ErrorResponse{Error: err.Error()})
		return true
	}
	return false
}
