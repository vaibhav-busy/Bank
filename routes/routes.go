package routes

import (
	"bank/handlers"

	"github.com/gin-gonic/gin"
)

func CreateRoutes() {
	router := gin.Default()

	bankRoute := router.Group("/bank")

	bankRoute.POST("", handlers.CreateBank)
	bankRoute.GET("", handlers.GetAllBanks)
	bankRoute.GET("/branch", handlers.GetAllBanksWithBranches)
	bankRoute.GET("/:id", handlers.GetBankByID)
	bankRoute.GET("/:id/branch", handlers.GetAllBranchesOfBankByID)
	bankRoute.PATCH("", handlers.UpdateBank)
	bankRoute.DELETE("", handlers.DeleteAllBanks)
	bankRoute.DELETE("/:id", handlers.DeleteBankByID)

	branchRoute := router.Group("/branch")

	branchRoute.POST("", handlers.CreateBranch)
	branchRoute.GET("", handlers.GetAllBranches)
	branchRoute.GET("/bank/account", handlers.GetAllBranchesWithBankAndAccounts)
	branchRoute.GET("/:id", handlers.GetBranchByID)
	branchRoute.GET("/:id/account", handlers.GetAllAccountsOfBranchByID)
	branchRoute.PATCH("", handlers.UpdateBranch)
	branchRoute.DELETE("", handlers.DeleteAllBranches)
	branchRoute.DELETE("/:id", handlers.DeleteBranchByID)

	// add joint account
	customerRoute := router.Group("/customer")

	customerRoute.POST("", handlers.CreateCustomer)
	customerRoute.GET("", handlers.GetAllCustomers)
	customerRoute.GET("/:id", handlers.GetCustomerByID)
	customerRoute.GET("/:id/account", handlers.GetAllAccountsByCustomerID)
	customerRoute.PATCH("", handlers.UpdateCustomer)
	customerRoute.DELETE("", handlers.DeleteAllCustomers)
	customerRoute.DELETE("/:id", handlers.DeleteCustomerByID)

	//get all customers and accounts
	accountRoute := router.Group("/account")

	accountRoute.POST("", handlers.CreateAccount)
	accountRoute.GET("", handlers.GetAllAccounts)
	accountRoute.GET("/:id", handlers.GetAccountByID)
	accountRoute.GET("/:id/customer", handlers.GetAllCustomersByAccountID)
	accountRoute.PATCH("", handlers.UpdateAccount)
	accountRoute.DELETE("", handlers.DeleteAllAccounts)
	accountRoute.DELETE("/:id", handlers.DeleteAccountByID)

	mappingRoute := router.Group("/account_to_customer")

	mappingRoute.POST("", handlers.CreateMapping)
	mappingRoute.GET("", handlers.GetAllMappings)
	mappingRoute.GET("/:id", handlers.GetMappingByID)
	mappingRoute.PATCH("", handlers.UpdateMapping)
	mappingRoute.DELETE("", handlers.DeleteAllMappings)
	mappingRoute.DELETE("/:id", handlers.DeleteMappingByID)

	transactionRoute := router.Group("/transaction")
	
	transactionRoute.POST("", handlers.CreateTransaction)
	transactionRoute.GET("", handlers.GetAllTransactions)
	transactionRoute.GET("/:id", handlers.GetTransactionByID)
	transactionRoute.GET("/account/:id", handlers.GetTransactionByAccountID)
	transactionRoute.DELETE("", handlers.DeleteAllTransactions)
	transactionRoute.DELETE("/:id", handlers.DeleteTransactionByID)

	router.Run(":8080")
}
