package storage_equipment

import (
	"context"
	"net/http"
	"time"

	"github.com/filipskrabak/storage-equipment-webapi/internal/db_service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type implEquipmentManagementAPI struct {
}

func NewEquipmentManagementApi() EquipmentManagementAPI {
	return &implEquipmentManagementAPI{}
}

func (o implEquipmentManagementAPI) CreateEquipment(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database service not available"})
		return
	}

	db, ok := value.(db_service.DbService[EquipmentItem])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database service type assertion failed"})
		return
	}

	var equipment EquipmentCreate
	if err := c.ShouldBindJSON(&equipment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	// Validate required fields
	if equipment.Name == "" || equipment.SerialNumber == "" || equipment.Manufacturer == "" || equipment.Location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields: name, serialNumber, manufacturer, location"})
		return
	}

	// Create new equipment item
	equipmentItem := EquipmentItem{
		Id:               uuid.New().String(),
		Name:             equipment.Name,
		SerialNumber:     equipment.SerialNumber,
		Manufacturer:     equipment.Manufacturer,
		Model:            equipment.Model,
		InstallationDate: equipment.InstallationDate,
		Location:         equipment.Location,
		ServiceInterval:  equipment.ServiceInterval,
		LastService:      equipment.LastService,
		LifeExpectancy:   equipment.LifeExpectancy,
		Status:           equipment.Status,
		Notes:            equipment.Notes,
	}

	// Set default status if not provided
	if equipmentItem.Status == "" {
		equipmentItem.Status = "operational"
	}

	// Calculate next service date if service interval and last service are provided
	if equipmentItem.ServiceInterval > 0 && equipmentItem.LastService != "" {
		if lastServiceTime, err := time.Parse("2006-01-02", equipmentItem.LastService); err == nil {
			nextService := lastServiceTime.AddDate(0, 0, int(equipmentItem.ServiceInterval))
			equipmentItem.NextService = nextService.Format("2006-01-02")
		}
	}

	if err := db.CreateDocument(context.Background(), equipmentItem.Id, &equipmentItem); err != nil {
		if err == db_service.ErrConflict {
			c.JSON(http.StatusConflict, gin.H{"error": "Equipment with this ID already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create equipment", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, equipmentItem)
}

func (o implEquipmentManagementAPI) DeleteEquipment(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database service not available"})
		return
	}

	db, ok := value.(db_service.DbService[EquipmentItem])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database service type assertion failed"})
		return
	}

	equipmentId := c.Param("equipmentId")
	if equipmentId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Equipment ID is required"})
		return
	}

	if err := db.DeleteDocument(context.Background(), equipmentId); err != nil {
		if err == db_service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete equipment", "details": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (o implEquipmentManagementAPI) GetAllEquipment(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database service not available"})
		return
	}

	db, ok := value.(db_service.DbService[EquipmentItem])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database service type assertion failed"})
		return
	}

	equipment, err := db.FindAllDocuments(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve equipment", "details": err.Error()})
		return
	}

	// Convert []*EquipmentItem to []EquipmentItem
	var equipmentList []EquipmentItem
	for _, item := range equipment {
		if item != nil {
			equipmentList = append(equipmentList, *item)
		}
	}

	c.JSON(http.StatusOK, equipmentList)
}

func (o implEquipmentManagementAPI) GetEquipmentById(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database service not available"})
		return
	}

	db, ok := value.(db_service.DbService[EquipmentItem])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database service type assertion failed"})
		return
	}

	equipmentId := c.Param("equipmentId")
	if equipmentId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Equipment ID is required"})
		return
	}

	equipment, err := db.FindDocument(context.Background(), equipmentId)
	if err != nil {
		if err == db_service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve equipment", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, equipment)
}

func (o implEquipmentManagementAPI) UpdateEquipment(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database service not available"})
		return
	}

	db, ok := value.(db_service.DbService[EquipmentItem])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database service type assertion failed"})
		return
	}

	equipmentId := c.Param("equipmentId")
	if equipmentId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Equipment ID is required"})
		return
	}

	// Get existing equipment
	existingEquipment, err := db.FindDocument(context.Background(), equipmentId)
	if err != nil {
		if err == db_service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Equipment not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve equipment", "details": err.Error()})
		}
		return
	}

	var updateData EquipmentUpdate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	// Update only provided fields
	if updateData.Name != "" {
		existingEquipment.Name = updateData.Name
	}
	if updateData.SerialNumber != "" {
		existingEquipment.SerialNumber = updateData.SerialNumber
	}
	if updateData.Manufacturer != "" {
		existingEquipment.Manufacturer = updateData.Manufacturer
	}
	if updateData.Model != "" {
		existingEquipment.Model = updateData.Model
	}
	if updateData.Location != "" {
		existingEquipment.Location = updateData.Location
	}
	if updateData.ServiceInterval > 0 {
		existingEquipment.ServiceInterval = updateData.ServiceInterval
	}
	if updateData.LastService != "" {
		existingEquipment.LastService = updateData.LastService

		// Recalculate next service if service interval is set
		if existingEquipment.ServiceInterval > 0 {
			if lastServiceTime, err := time.Parse("2006-01-02", existingEquipment.LastService); err == nil {
				nextService := lastServiceTime.AddDate(0, 0, int(existingEquipment.ServiceInterval))
				existingEquipment.NextService = nextService.Format("2006-01-02")
			}
		}
	}
	if updateData.LifeExpectancy > 0 {
		existingEquipment.LifeExpectancy = updateData.LifeExpectancy
	}
	if updateData.Status != "" {
		existingEquipment.Status = updateData.Status
	}
	if updateData.Notes != "" {
		existingEquipment.Notes = updateData.Notes
	}

	if err := db.UpdateDocument(context.Background(), equipmentId, existingEquipment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update equipment", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingEquipment)
}
