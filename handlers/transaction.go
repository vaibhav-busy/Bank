package handlers

import (
	// "fmt"
	// "fmt"
	"errors"
	"net/http"
	"strconv"

	"bank/database"
	"bank/models"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

func CreateTransaction(c *gin.Context) {

	tx, txErr := database.Db.Begin()

	if txErr != nil {
		// tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't create transaction"})
		return
	}

	var transaction models.Transaction
	if err := c.ShouldBind(&transaction); err != nil {
		// tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	_, insertErr := database.Db.Model(&transaction).Returning("id").Insert()

	if insertErr != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": insertErr.Error()})
		return
	}
	// fmt.Print(transaction.AccountID)
	var newBal float64
	var err error

	switch transaction.Mode {
	case "deposit":
		newBal, err = Deposit(tx, transaction.AccountID, transaction.Amount)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": insertErr.Error()})
			return
		}
		tx.Commit()
		c.JSON(http.StatusCreated, gin.H{"Transaction created successfully and new balance is : ": newBal})
	case "withdraw":
		newBal, err = Withdraw(tx, transaction.AccountID, transaction.Amount)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tx.Commit()
		c.JSON(http.StatusCreated, gin.H{"Transaction created successfully and new balance is : ": newBal})

	case "transfer":

		if transaction.ReceiverAccountNo == 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("receiver account no is not sent").Error()})
			return
		}

		err = Transfer(tx, transaction.AccountID, transaction.ReceiverAccountNo, transaction.Amount)

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tx.Commit()

		c.JSON(http.StatusCreated, gin.H{"Transaction completed successfully": ""})
	}
}

func GetAllTransactions(c *gin.Context) {
	var transactions []models.Transaction
	getErr := database.Db.Model(&transactions).Select()

	if getErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": getErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Transactions": transactions})
}

func GetTransactionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	transaction := new(models.Transaction)
	tranErr := database.Db.Model(transaction).Where("id = ?", id).Select()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": tranErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Transaction": transaction})
}

func GetTransactionByAccountID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	transaction := new([]models.Transaction)
	tranErr := database.Db.Model(transaction).Where("account_id = ?", id).Select()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": tranErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Transaction": *transaction})
}

func DeleteAllTransactions(c *gin.Context) {
	transaction := new(models.Transaction)

	res, err := database.Db.Query(transaction, "DELETE from transactions")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total rows deleted ": res.RowsAffected()})
}

func DeleteTransactionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	transaction := new(models.Transaction)

	_, delErr := database.Db.Model(transaction).Where("id = ?", id).Returning("id").Delete()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": delErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Deleted Transaction with id = ": transaction.ID})
}

func Deposit(tx *pg.Tx, id uint64, amount float64) (float64, error) {
	account := new(models.Account)
	account.ID = id
	err := tx.Model(account).WherePK().Select()
	if err != nil {
		return account.Balance, err
	}

	account.Balance += amount
	_, err = tx.Model(account).Column("balance").WherePK().Returning("balance").Update()
	if err != nil {
		return account.Balance, err
	}
	return account.Balance, nil
}

func Withdraw(tx *pg.Tx, id uint64, amount float64) (float64, error) {
	account := new(models.Account)
	account.ID = id
	err := tx.Model(account).WherePK().Select()
	if err != nil {
		return account.Balance, err
	}
	if account.Balance < amount {
		return account.Balance, errors.New("account balance is less than amount")
	}
	account.Balance -= amount
	_, err = tx.Model(account).Column("balance").WherePK().Returning("balance").Update()
	if err != nil {
		return account.Balance, err
	}
	return account.Balance, nil
}

func Transfer(tx *pg.Tx, senderId uint64, receiverAccNo uint64, amount float64) error {

	var account models.Account

	res, err := tx.Query(account, "UPDATE accounts SET balance = balance - ? WHERE id = ? AND balance - ?0 >= 0", amount, senderId)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("not a valid request (either id is wrong or balance is not sufficient)")
	}

	_, err = tx.Query(account, "UPDATE accounts SET balance = balance + ? WHERE account_no = ?", amount, receiverAccNo)
	if err != nil {
		return err
	}
	return nil
}
