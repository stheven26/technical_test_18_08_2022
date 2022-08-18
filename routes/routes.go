package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stheven26/controllers"
)

func Router(app *fiber.App) {
	app.Post("/api/v1/register", controllers.Register)
	app.Post("/api/v1/login", controllers.Login)
	app.Post("/api/v1/logout", controllers.Logout)
	app.Get("/api/v1/blogs/", controllers.GetAllBlog)
	app.Get("/api/v1/blogs/:id", controllers.GetBlogById)
	app.Post("/api/v1/blogs/", controllers.CreateBlog)
	app.Put("/api/v1/blogs/:id", controllers.UpdateBlog)
	app.Delete("/api/v1/blogs/:id", controllers.DeleteBlog)
}
