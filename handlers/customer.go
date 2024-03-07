package handlers

import (
	"net/http"
	"strconv"

	"bank/database"
	"bank/models"

	"github.com/gin-gonic/gin"
)

func CreateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBind(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	_, insertErr := database.Db.Model(&customer).Returning("id").Insert()

	if insertErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": insertErr.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"Customer created with id : ": customer.ID})
}

func GetAllCustomers(c *gin.Context) {

	var customers []models.Customer
	getErr := database.Db.Model(&customers).Select()

	if getErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": getErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Customers": customers})
}

func GetAllAccountsByCustomerID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	var accounts []models.Account
	customer_accounts := database.Db.Model((*models.AccountToCustomer)(nil)).ColumnExpr("account_id").Where("customer_id = ?", id)
	getErr := database.Db.Model(&accounts).Where("id IN (?)", customer_accounts).Select()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": getErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Accounts of customer ": accounts})
}

func GetCustomerByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	customer := new(models.Customer)
	getErr := database.Db.Model(customer).Where("id = ?", id).Select()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": getErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Customer": *customer})
}

func UpdateCustomer(c *gin.Context) {

	customer := new(models.Customer)
	err := c.ShouldBindJSON(customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, txErr := database.Db.Begin()

	if txErr != nil {
		return
	}

	res, err := tx.Model(customer).WherePK().Returning("id").Update()
	if err != nil {
		tx.Rollback()
		return
	}
	if res.RowsAffected() == 0 {
		tx.Rollback()
		return
	}
	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{"Update customer with id = ": customer.ID})
}

func DeleteAllCustomers(c *gin.Context) {
	customer := new(models.Customer)

	res, err := database.Db.Query(customer, "DELETE from customers")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total rows deleted ": res.RowsAffected()})
}

func DeleteCustomerByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	customer := new(models.Customer)

	_, delErr := database.Db.Model(customer).Where("id = ?", id).Returning("id").Delete()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": delErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Deleted Customer with id = ": customer.ID})
}
