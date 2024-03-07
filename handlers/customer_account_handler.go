package handlers

import (
	"net/http"
	"strconv"

	"bank/database"
	"bank/models"

	"github.com/gin-gonic/gin"
)

func CreateMapping(c *gin.Context) {
	var mapping models.AccountToCustomer

	if err := c.ShouldBind(&mapping); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	_, insertErr := database.Db.Model(&mapping).Returning("id").Insert()

	if insertErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": insertErr.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Mapping created with id : ": mapping.ID})
}

func GetAllMappings(c *gin.Context) {
	var mappings []models.AccountToCustomer
	getErr := database.Db.Model(&mappings).Select()

	if getErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": getErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Mappings": mappings})
}

func GetMappingByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	mapping := new(models.AccountToCustomer)
	getErr := database.Db.Model(mapping).Where("id = ?", id).Select()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": getErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Mapping": *mapping})
}

func UpdateMapping(c *gin.Context) {
	mapping := new(models.AccountToCustomer)
	err := c.ShouldBindJSON(mapping)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, txErr := database.Db.Begin()
	if txErr != nil {
		return 
	}

	res, err := tx.Model(mapping).WherePK().Returning("id").Update()
	if err != nil {
		tx.Rollback()
		return
	}
	if res.RowsAffected() == 0 {
		tx.Rollback()
		return
	}
	tx.Commit()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Update mapping with id = ": mapping.ID})
}

func DeleteAllMappings(c *gin.Context) {
	mapping := new(models.AccountToCustomer)

	res, err := database.Db.Query(mapping, "DELETE from mappings")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total rows deleted ": res.RowsAffected()})
}

func DeleteMappingByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	mapping := new(models.AccountToCustomer)

	_, delErr := database.Db.Model(mapping).Where("id = ?", id).Returning("id").Delete()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": delErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Deleted Mapping with id = ": mapping.ID})
}
