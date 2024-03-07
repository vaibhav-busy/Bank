package handlers

import (
	"bank/database"
	"bank/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)


func CreateBranch(c *gin.Context) {
	var branch models.Branch

	if err := c.ShouldBind(&branch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	_, insertErr := database.Db.Model(&branch).Returning("id").Insert()

	if insertErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": insertErr.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Branch created with id : ": &branch.ID})

}


func GetAllBranches(c *gin.Context) {
	var branches []models.Branch
	getErr := database.Db.Model(&branches).Select()

	if getErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": getErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"branch": branches})
}


func GetAllBranchesWithBankAndAccounts(c *gin.Context) {
	var branches []models.Branch
	getErr := database.Db.Model(&branches).Relation("Accounts").Relation("Bank").Select()


	if getErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": getErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"branch": branches})
}

func GetBranchByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64) // 10 is for base 10 digits and 64 for uint64  bit size
	if err != nil {
		panic(err)
	}
	branch := new(models.Branch)
	getErr := database.Db.Model(branch).Where("id = ?", id).Select()
	if getErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": getErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"branch": branch})
}

func GetAllAccountsOfBranchByID(c *gin.Context){
	id, err := strconv.ParseUint(c.Param("id"), 10, 64) // 10 is for base 10 digits and 64 for uint64  bit size
	if err != nil {
		panic(err)
	}
	accounts := new([]models.Account)
	getErr := database.Db.Model(accounts).Where("branch_id = ?", id).Select()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": getErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

func UpdateBranch(c *gin.Context) {

	branch := new(models.Branch)
	err := c.ShouldBindJSON(branch)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, txErr := database.Db.Begin()

	if txErr != nil {
		return 
	}

	res, err := tx.Model(branch).WherePK().Returning("id").Update()
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
	c.JSON(http.StatusCreated, gin.H{"Update branch with id = ": branch.ID})
}


func DeleteAllBranches(c *gin.Context) {
	branch := new(models.Branch)


	res, err := database.Db.Query(branch, "DELETE from branches")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total rows deleted ": res.RowsAffected()})
}


func DeleteBranchByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	branch := new(models.Branch)

	_, delErr := database.Db.Model(branch).Where("id = ?", id).Returning("id").Delete()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": delErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Deleted branch with id = ": branch.ID})
}