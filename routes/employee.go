package routes

import (
	"github.com/gorilla/mux"
	"github.com/sabarish-manoharan/emp-management/controllers"
	"github.com/sabarish-manoharan/emp-management/middleware"
)

func RegisterUserRoutes(router *mux.Router) {
	router.HandleFunc("/register", controllers.RegisterUser).Methods("POST");
}

func LoginUserRoutes(router *mux.Router){
	router.HandleFunc("/login",controllers.LoginUser).Methods("POST");
}
func GetUserRoutes(router *mux.Router){
	router.HandleFunc("/users",controllers.GetUser).Methods("GET");
}



func RegisterEmployeeRoutes(r *mux.Router) {
	secured := r.PathPrefix("/api").Subrouter()
	secured.Use(middleware.AuthMiddleware);

	secured.HandleFunc("/employees",  controllers.GetEmployee).Methods("GET")
	secured.HandleFunc("/employee",  controllers.CreateEmployee).Methods("POST")
	secured.HandleFunc("/employee/{id}",controllers.UpdateEmployee).Methods("PUT")
	secured.HandleFunc("/employee/{id}",controllers.DeleteEmployee).Methods("DELETE")
}
