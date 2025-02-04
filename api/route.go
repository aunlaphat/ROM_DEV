package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "boilerplate-backend-go/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

// Define route of service and middleware relate
func (app *Application) routes() http.Handler {
	router := chi.NewRouter()
	//CORS Middleware
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Allow only localhost for testing and production ip address
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},

		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum cache age (in seconds)
	})
	router.Use(cors.Handler)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	apiRouter := chi.NewRouter()
	apiRouter.Get("/", app.test)
	apiRouter.Get("/swagger/*", httpSwagger.Handler())
	// Auth router
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "uploads"))
	FileServer(apiRouter, "/uploads", filesDir)
	app.AuthRoute(apiRouter)
	// app.UserTestRoute(apiRouter)
	app.FileServerRoute(apiRouter)
	//app.Constants(apiRouter)
	app.Excels(apiRouter)
	app.UserRoute(apiRouter)
	app.BeforeReturnRoute(apiRouter)
	app.ReturnOrders(apiRouter)
	app.ImportOrderRoute(apiRouter)
	app.TradeReturnRoute(apiRouter)

	router.Mount("/api", apiRouter)
	return router

	//Need Auth Group
	// apiRouter.Group(func(router chi.Router) {
	// 	router.Use(jwtauth.Verifier(app.TokenAuth))
	// 	router.Use(jwtauth.Authenticator)
	// 	//Router Example
	// 	router.Route("/wishlist", func(r chi.Router) {
	// 	//Partial Middeleware Example

	// 		r.Get("/", app.GetWishlist) //This route any department can access
	// 		r.Group(func(r chi.Router) {
	// 			r.Use(requireDepartment("G01")) //Below route only G01 department can access
	// 			r.Post("/", app.CreateWishlist)
	// 			r.Get("/{id}", app.GetWishlistByID)
	// 			r.Patch("/{id}", app.UpdateWishlistByID)
	// 	})

	// 	//Router Example 2

	// 		router.Route("/prepare_pr", func(r chi.Router) {
	// 			r.Use(requireDepartment("G01")) Require valid department to access
	// 			r.Get("/", app.GetPreparePRList)
	// 			r.Post("/", app.InsertPreparePRLine)
	// 			r.Patch("/", app.UpdatePreparePR)
	// 			r.Delete("/", app.DeletePreparePR)
	// 			r.Post("/{id}", app.CreatePreparePR)
	// 		})

	// })
}

func (app *Application) test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the API!"))
}
