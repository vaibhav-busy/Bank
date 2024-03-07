package handlers

import (
	"net/http"
	"strconv"

	"bank/models"
    "bank/database"
	"github.com/gin-gonic/gin"
)

func CreateBank(c *gin.Context) {
	var bank models.Bank
	if err := c.ShouldBind(&bank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Insert new bank and get the ID
	_, insertErr := database.Db.Model(&bank).Returning("id").Insert()
	if insertErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": insertErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Bank created with id": bank.ID})
}

func GetAllBanks(c *gin.Context) {

	var banks []models.Bank
	err := database.Db.Model(&banks).Select()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Banks": banks})
}

func GetAllBanksWithBranches(c *gin.Context) {
	
	var banks []models.Bank
	err := database.Db.Model(&banks).Relation("Branches").Select()


	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Banks": banks})
}

func GetBankByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	bank := new(models.Bank)
	getErr := database.Db.Model(bank).Where("id = ?", id).Select()

	if getErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": getErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Bank": bank})
}

func GetAllBranchesOfBankByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	branches := new([]models.Branch)
	getErr := database.Db.Model(branches).Where("bank_id = ?", id).Select()

	if getErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": getErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Branches : ": branches})
}

func UpdateBank(c *gin.Context) {
	bank := new(models.Bank)
	err := c.ShouldBindJSON(bank)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// updated_bank_id, err := queries.UpdateBank(bank)
	tx, txErr := database.Db.Begin()
	if txErr != nil {
		return 
	}

	res, err := tx.Model(bank).WherePK().Returning("id").Update()
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
	c.JSON(http.StatusCreated, gin.H{"Update bank with id = ":bank.ID})
}

func DeleteAllBanks(c *gin.Context) {
	bank := new(models.Bank)

	res, err := database.Db.Query(bank, "DELETE from banks")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total rows deleted ": res.RowsAffected()})
}

func DeleteBankByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	// deleted_id, err := queries.DeleteBankByID(id)
	bank := new(models.Bank)

	_, delErr := database.Db.Model(bank).Where("id = ?", id).Returning("id").Delete()
	if delErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": delErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Deleted Bank with id = ": bank.ID})
}
