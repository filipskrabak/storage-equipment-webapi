package storage_equipment

import (
	"context"
	"net/http"
	"time"

	"github.com/filipskrabak/storage-equipment-webapi/internal/db_service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type implEquipmentOrdersAPI struct {
}

func NewEquipmentOrdersApi() EquipmentOrdersAPI {
	return &implEquipmentOrdersAPI{}
}

func (o implEquipmentOrdersAPI) CreateOrder(c *gin.Context) {
	value, exists := c.Get("orders_db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Orders database service not available"})
		return
	}

	db, ok := value.(db_service.DbService[Order])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Orders database service type assertion failed"})
		return
	}

	var orderCreate OrderCreate
	if err := c.ShouldBindJSON(&orderCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	// Validate required fields
	if orderCreate.RequestedBy == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RequestedBy is required"})
		return
	}

	if len(orderCreate.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one item is required"})
		return
	}

	// Create new order
	order := Order{
		Id:                  uuid.New().String(),
		Items:               orderCreate.Items,
		RequestedBy:         orderCreate.RequestedBy,
		RequestorDepartment: orderCreate.RequestorDepartment,
		Status:              "pending",
		Notes:               orderCreate.Notes,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := db.CreateDocument(context.Background(), order.Id, &order); err != nil {
		if err == db_service.ErrConflict {
			c.JSON(http.StatusConflict, gin.H{"error": "Order with this ID already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (o implEquipmentOrdersAPI) ListOrders(c *gin.Context) {
	value, exists := c.Get("orders_db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Orders database service not available"})
		return
	}

	db, ok := value.(db_service.DbService[Order])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Orders database service type assertion failed"})
		return
	}

	orders, err := db.FindAllDocuments(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders", "details": err.Error()})
		return
	}

	var ordersList []Order
	for _, order := range orders {
		if order != nil {
			ordersList = append(ordersList, *order)
		}
	}

	c.JSON(http.StatusOK, ordersList)
}

func (o implEquipmentOrdersAPI) GetOrderById(c *gin.Context) {
	value, exists := c.Get("orders_db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Orders database service not available"})
		return
	}

	db, ok := value.(db_service.DbService[Order])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Orders database service type assertion failed"})
		return
	}

	orderId := c.Param("orderId")
	if orderId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	order, err := db.FindDocument(context.Background(), orderId)
	if err != nil {
		if err == db_service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, order)
}

func (o implEquipmentOrdersAPI) UpdateStatus(c *gin.Context) {
	value, exists := c.Get("orders_db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Orders database service not available"})
		return
	}

	db, ok := value.(db_service.DbService[Order])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Orders database service type assertion failed"})
		return
	}

	orderId := c.Param("orderId")
	if orderId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	// Get existing order
	existingOrder, err := db.FindDocument(context.Background(), orderId)
	if err != nil {
		if err == db_service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order", "details": err.Error()})
		}
		return
	}

	var updateData OrderUpdate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"pending":   true,
		"delivered": true,
		"cancelled": true,
	}

	if !validStatuses[updateData.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be: pending, delivered, or cancelled"})
		return
	}

	// Update status and timestamp
	if updateData.RequestedBy != "" {
		existingOrder.RequestedBy = updateData.RequestedBy
	}
	if updateData.RequestorDepartment != "" {
		existingOrder.RequestorDepartment = updateData.RequestorDepartment
	}
	if updateData.Notes != "" {
		existingOrder.Notes = updateData.Notes
	}
	if updateData.Items != nil {
		existingOrder.Items = updateData.Items
	}

	existingOrder.Status = updateData.Status
	existingOrder.UpdatedAt = time.Now()

	if err := db.UpdateDocument(context.Background(), orderId, existingOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingOrder)
}

func (o implEquipmentOrdersAPI) CancelOrder(c *gin.Context) {
	value, exists := c.Get("orders_db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Orders database service not available"})
		return
	}

	db, ok := value.(db_service.DbService[Order])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Orders database service type assertion failed"})
		return
	}

	orderId := c.Param("orderId")
	if orderId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	if err := db.DeleteDocument(context.Background(), orderId); err != nil {
		if err == db_service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order", "details": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
