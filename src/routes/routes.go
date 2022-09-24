package routes

import (
	"ambassador/src/controllers"
	"ambassador/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")

	admin := api.Group("admin")

	admin.Post("/register", controllers.Register)
	admin.Post("/login", controllers.Login)

	adminAuthenticated := admin.Use(middlewares.IsAuthenticated)
	adminAuthenticated.Get("/me", controllers.User)
	adminAuthenticated.Post("/logout", controllers.Logout)
	adminAuthenticated.Put("/me", controllers.UpdateInfo)
	adminAuthenticated.Put("/me/password", controllers.UpdateInfo)
	adminAuthenticated.Get("/ambassadors", controllers.Ambassador)
	adminAuthenticated.Get("/products", controllers.Products)
	adminAuthenticated.Post("/products", controllers.CreateProduct)
	adminAuthenticated.Get("/products/:id", controllers.GetProduct)
	adminAuthenticated.Put("/products/:id", controllers.UpdateProduct)
	adminAuthenticated.Delete("/products/:id", controllers.DeleteProduct)
	adminAuthenticated.Get("/users/:id/links", controllers.Link)
	adminAuthenticated.Get("/orders", controllers.Orders)

	ambassador := api.Group("ambassador")

	ambassador.Post("/register", controllers.Register)
	ambassador.Post("/login", controllers.Login)

	ambassadorAuthenticated := ambassador.Use(middlewares.IsAuthenticated)
	ambassadorAuthenticated.Get("/me", controllers.User)
	ambassadorAuthenticated.Post("/logout", controllers.Logout)
	ambassadorAuthenticated.Put("/me", controllers.UpdateInfo)
	ambassadorAuthenticated.Put("/me/password", controllers.UpdateInfo)
}
