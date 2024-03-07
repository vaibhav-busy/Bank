package handlers
												//update tranfer,withdraw functions
import (
	"bank/database"
	"bank/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateAccount(c *gin.Context) {
	var account models.Account
	if err := c.ShouldBind(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	_, insertErr := database.Db.Model(&account).Returning("id").Insert()
	if insertErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": insertErr.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Account created with id : ": &account.ID})
}

func GetAllAccounts(c *gin.Context) {
	var accounts []models.Account
	getErr := database.Db.Model(&accounts).Select()


	if getErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": getErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Accounts": accounts})
}

func GetAllCustomersByAccountID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	var customers []models.Customer
	customer_accounts := database.Db.Model((*models.AccountToCustomer)(nil)).ColumnExpr("customer_id").Where("account_id = ?",id)
	getErr := database.Db.Model(&customers).Where("id IN (?)",customer_accounts).Select()

	if getErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": getErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Customers of accounts ": customers})
}

func GetAccountByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	account := new(models.Account)
	getErr := database.Db.Model(account).Where("id = ?", id).Select()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": getErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Account": *account})
}

func UpdateAccount(c *gin.Context) {
	account := new(models.Account)
	err := c.ShouldBindJSON(account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, txErr := database.Db.Begin()

	if txErr != nil {
		return 
	}

	res, err := tx.Model(account).WherePK().Returning("id").Update()
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
	c.JSON(http.StatusCreated, gin.H{"Update account with id = ": account.ID})
}

func DeleteAllAccounts(c *gin.Context) {
	account := new(models.Account)

	res, err := database.Db.Query(account, "DELETE from accounts")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total rows deleted ": res.RowsAffected()})
}

func DeleteAccountByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	account := new(models.Account)

	_, delErr := database.Db.Model(account).Where("id = ?", id).Returning("id").Delete()

	if delErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": delErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Deleted Account with id = ": account.ID})
}
