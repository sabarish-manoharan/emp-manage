package controllers

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sabarish-manoharan/emp-management/db"
	"github.com/sabarish-manoharan/emp-management/models"
)

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var emp models.Employee

	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !dataValidation(w, emp) {
		http.Error(w, "Invalid Data", http.StatusBadRequest)
		return
	}
	if err := db.DB.Create(&emp).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondJSON(w, 200, emp)
}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	var employees []models.Employee
	if err := db.DB.Find(&employees).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, 200, employees)
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var emp models.Employee

	if err := db.DB.First(&emp, id).Error; err != nil {
		http.Error(w, "Employee Not Found", http.StatusNotFound)
		return
	}

	var updatedEmployee models.Employee
	if err := json.NewDecoder(r.Body).Decode(&updatedEmployee); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !dataValidation(w, updatedEmployee) {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	emp.Name = updatedEmployee.Name
	emp.Age = updatedEmployee.Age
	emp.Email = updatedEmployee.Email
	emp.Role = updatedEmployee.Role

	if err := db.DB.Save(&emp).Error; err != nil {
		http.Error(w, "Failed to Update employee", http.StatusInternalServerError)
		return
	}
	respondJSON(w, 200, emp)
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := db.DB.Delete(&models.Employee{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, 200, map[string]interface{}{
		"success": true,
		"message": "Employee Deleted Successfully",
	})
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
func dataValidation(w http.ResponseWriter, emp models.Employee) bool {

	if isEmptyOrWhiteSpaces(emp.Name) || isEmptyOrWhiteSpaces(emp.Role) || isEmptyOrWhiteSpaces(emp.Email) {
		http.Error(w, "Missing Field Required", http.StatusBadRequest)
		return false
	}
	if !validEmail(emp.Email) {
		http.Error(w, "Invalid Email", http.StatusBadRequest)
		return false
	}
	if emp.Age < 1 || emp.Age > 100 {
		http.Error(w, "Age must be above 1 and below 100", http.StatusBadRequest)
		return false
	}
	return true
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}

func isEmptyOrWhiteSpaces(s string) bool {
	return strings.TrimSpace(s) == ""
}
