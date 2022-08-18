package server

import (
	"audioPhile/handlers"
	"audioPhile/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Server struct {
	chi.Router
}

func SetupRoutes() *Server {
	router := chi.NewRouter()
	router.Route("/api", func(api chi.Router) {
		api.Post("/signup", handlers.SignUp)
		api.Post("/login", handlers.Login)
		api.Route("/audiophile", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)
			r.Route("/admin", func(admin chi.Router) {
				admin.Use(middleware.AdminMiddleware)
				admin.Get("/products", handlers.GetAllProducts)
				admin.Route("/discount", func(coupon chi.Router) {
					coupon.Post("/addCoupon", handlers.CreateCoupon)

				})
				admin.Post("/category", handlers.CreateProductCategory)
				admin.Route("/{categoryID}", func(category chi.Router) {
					category.Put("/", handlers.UpdateProductCategory)
					category.Delete("/", handlers.DeleteProductCategory)
					category.Post("/inventory", handlers.CreateProductInventory)
					category.Route("/{inventoryID}", func(inventory chi.Router) {
						inventory.Put("/", handlers.UpdateProductInventory)
						inventory.Delete("/", handlers.DeleteProductInventory)
						inventory.Post("/", handlers.CreateProduct)
						inventory.Route("/{productID}", func(product chi.Router) {
							product.Post("/", handlers.UploadProductImage)
							product.Get("/", handlers.GetAllImageDetails)
							product.Delete("/", handlers.DeleteProduct)
							product.Put("/", handlers.UpdateProduct)

						})
					})
				})
			})
		})
		api.Route("/user", func(user chi.Router) {
			user.Get("/products", handlers.GetAllProducts)
			user.Route("/{productID}", func(image chi.Router) {
				image.Get("/", handlers.GetAllImageDetails)
			})
			user.Route("/cart", func(cart chi.Router) {
				cart.Use(middleware.AuthMiddleware)
				cart.Route("/address", func(address chi.Router) {
					address.Post("/", handlers.AddUserAddress)
					address.Route("/{addressID}", func(add chi.Router) {
						add.Put("/", handlers.UpdateUserAddress)
						add.Delete("/", handlers.DeleteUserAddress)
					})
				})

				cart.Post("/addProduct", handlers.AddProductToCart)
				cart.Route("/{productID}", func(new chi.Router) {
					//new.Use(middleware.UserAddressAndQuantityValidateMiddleware)
					new.Route("/buy", func(order chi.Router) {
						order.Put("/", handlers.BuyProduct)
						order.Post("/", handlers.AddPaymentDetails)
						order.Route("/{paymentID}", func(payment chi.Router) {
							payment.Get("/", handlers.AddOrderDetails)
						})
					})

				})
				cart.Get("/logout", handlers.SignOut)
			})

		})

	})
	return &Server{router}

}
func (svc *Server) Run(port string) error {
	return http.ListenAndServe(port, svc)
}
