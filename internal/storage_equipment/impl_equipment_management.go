package storage_equipment

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implEquipmentManagementAPI struct {
}

func NewEquipmentManagementApi() EquipmentManagementAPI {
	return &implEquipmentManagementAPI{}
}

func (o implEquipmentManagementAPI) CreateEquipment(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implEquipmentManagementAPI) DeleteEquipment(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implEquipmentManagementAPI) GetAllEquipment(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implEquipmentManagementAPI) GetEquipmentById(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implEquipmentManagementAPI) UpdateEquipment(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
