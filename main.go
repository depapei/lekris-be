package main

import (
	"Lekris-BE/controller"
	"Lekris-BE/controller/product"
	"Lekris-BE/controller/supply"
	"Lekris-BE/controller/transaction"
	"Lekris-BE/middleware"
	"Lekris-BE/model"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func main() {

	router := gin.New()

	router.Use(Cors())
	router.Use(gin.Logger())

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	model.ConnectDatabase()

	generateModel()
	api := router.Group("/api")
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/login/", controller.Login)
	}

	protectedRoutes := api
	protectedRoutes.Use(middleware.AuthMiddleware())
	{
		productRoutes := protectedRoutes.Group("/products")
		{
			productRoutes.GET("/", product.Index)
			productRoutes.GET("/:id/", product.Detail)
			productRoutes.POST("/", product.Create)
			productRoutes.PUT("/:id/", product.Update)
			productRoutes.DELETE("/:id/", product.Delete)
		}
		supplyRoutes := protectedRoutes.Group("/supplies")
		{
			supplyRoutes.GET("/", supply.Index)
			supplyRoutes.GET("/:id/", supply.Detail)
			supplyRoutes.POST("/", supply.Create)
			supplyRoutes.PUT("/:id/", supply.Update)
			supplyRoutes.DELETE("/:id/", supply.Delete)
		}
		transactionRoutes := protectedRoutes.Group("/transactions")
		{
			transactionRoutes.GET("/", transaction.Index)
			transactionRoutes.GET("/:id/", transaction.Detail)
			transactionRoutes.POST("/", transaction.Create)
			transactionRoutes.PUT("/:id/", transaction.Update)
			transactionRoutes.DELETE("/:id/", transaction.Delete)
		}
	}

	port := os.Getenv("PORT")
	router.Run(port)
}

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Timezone", "User", "X-Telegram-Auth-Date", "X-Telegram-Hash", "X-Telegram-Init-Data", "Service-Token", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Origin", "Cache-Control", "X-Requested-With"},
		AllowCredentials: false,
		ExposeHeaders:    []string{"Total-records", "Content-disposition"},
	})
}

func generateModel() {
	// Generate Models
	dsn := "host=localhost user=postgres password=db123 dbname=lele-krispy port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Initialize the generator
	g := gen.NewGenerator(gen.Config{
		OutPath: "./models",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(db)

	// Generate structs from all tables of the current database
	// The generated models will be placed in the specified OutPath
	user := g.GenerateModel("user")
	supply := g.GenerateModel("supply")
	products := g.GenerateModel("products")
	detail_transaction := g.GenerateModel("detail_transaction",
		gen.FieldRelate(field.HasOne, "Product", products, &field.RelateConfig{
			GORMTag: field.GormTag{"foreignKey": []string{"product_id"}},
		}),
	)
	transaction := g.GenerateModel("transactions",
		gen.FieldRelate(field.HasMany, "DetailTransaction", detail_transaction, &field.RelateConfig{
			GORMTag: field.GormTag{"foreignKey": []string{"transaction_id"}},
		}),
		gen.FieldType("isreturningcustomer", "*bool"),
	)
	g.ApplyBasic(user, supply, products, detail_transaction, transaction)

	g.GenerateAllTable()

	// Execute the generation
	g.Execute()
}
