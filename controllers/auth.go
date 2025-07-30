package controllers

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sabarish-manoharan/emp-management/db"
	"github.com/sabarish-manoharan/emp-management/models"
	//"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}
// var jwtSecret = []byte(viper.Get("JWT_SECRET_KEY").(string))
var JwtSecret = []byte("snow@10/mms");
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	validateUser(w, user)

	// 	//The cost determines how much computation is required to hash the password.
	// The higher the cost:
	// 🔒 More secure (harder to brute-force)
	// 🐢 Slower to compute (on purpose)

	// 	It is the default recommended cost (usually 10) set by the Go bcrypt package.
	// It provides a good balance between security and performance.
	// So:
	// For most apps: bcrypt.DefaultCost is safe and sufficient
	// For high-security systems: you might increase it (e.g. 12, 14), but test performanc

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost) // hashing password

	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword);
	
	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJson(w, 200, "Registered Successfully")
}

func LoginUser(w http.ResponseWriter, r *http.Request) {

	var LoginUser models.User

	if err := json.NewDecoder(r.Body).Decode(&LoginUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", LoginUser.Email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
        
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(LoginUser.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	jwt,err:=CreateJWT(user.Email,user.ID);
	if err !=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError);
		return
	}
	response := LoginResponse{
	Token:    jwt,
	Username: user.Name,
}

	respondJson(w, 200,response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	if err := db.DB.Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJson(w, 200, users)
}

func validateUser(w http.ResponseWriter, user models.User) {
	if CheckEmptyOrWhiteSpaces(user.Name) || CheckEmptyOrWhiteSpaces(user.Email) || CheckEmptyOrWhiteSpaces(user.Password) {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if !checkEmail(user.Email) {
		http.Error(w, "Invalid Email", http.StatusBadRequest)
		return
	}
	if len(user.Password) < 8 {
		http.Error(w, "Password is short", http.StatusBadRequest)
		return
	}
}

func checkEmail(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}
func CheckEmptyOrWhiteSpaces(s string) bool {
	return strings.TrimSpace(s) == ""
}

func respondJson(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func CreateJWT(email string, id uint) (string, error) {

	// 1️⃣ Claims: user data inside the token

	claims := jwt.MapClaims{
		"email": email,
		"id":    id,
		"exp":   time.Now().Add(time.Minute * 30).Unix(),
	}
     
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
	signedToken ,err := token.SignedString(JwtSecret);
	if err!=nil{
		return  "",err;
	}
	return signedToken,nil
}
