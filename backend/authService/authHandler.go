package authService

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cs_student_uni/microservice_store/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int32     `json:"user_id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedOn time.Time `json:"created_on"`
	Role      string    `json:"role"`
}

type Message struct {
	Status  int64  `json:"status"`
	Message string `json:"message"`
}

type Response struct {
	Status int64  `json:"status"`
	Data   []User `json:"data"`
}

type LoginResponse struct {
	Status   int64  `json:"status"`
	Token    string `json:"token"`
	Role     string `json:"role"`
	Username string `json:"username"`
}

var db *sql.DB

var jwtKey = []byte("eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiQWRtaW4iLCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6InN0dWRlbnRfY3MiLCJleHAiOjE3MzczMTcyMjksImlhdCI6MTcwNTY5NDgyOX0.U9ahjosw9nIhDY2c3NBmgjf4hB89RUzLG_ROIrtIEnM")

func Init() {
	db = database.GetDB()
	// http.HandleFunc("/login", login)
	// http.HandleFunc("/signup", signup)
	// http.HandleFunc("/users", getUsers)
	// http.HandleFunc("/authorize", authorize)
}

func saveUserToDB(username string, hashedPassword string, email string, created_on time.Time, role string) error {
	_, err := db.Exec("INSERT INTO users (username, password, email, created_on, role) VALUES ($1, $2, $3, $4, $5)", username, hashedPassword, email, time.Now(), role)
	if err != nil {
		fmt.Println("Error message: ", err)
		return err
	}
	return nil
}

func Signup(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.JSON(err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(err)
	}

	if user.Role == "" {
		user.Role = "customer" // Set the default role to "customer"
	}

	err = saveUserToDB(user.Username, string(hashedPassword), user.Email, user.CreatedOn, user.Role)
	if err != nil {
		return c.JSON("Error saving user to database")
	}

	return c.JSON(fiber.Map{"data": user})
}

// Assuming you have already established a connection to the PostgreSQL database and stored it in the variable db

func getUserPassword(username string) (string, error) {
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&hashedPassword)
	if err != nil {
		// Handle the error, such as returning a default value or logging the error
		return "", err
	}
	return hashedPassword, nil
}

func getUserRole(username string) (string, error) {
	var role string
	err := db.QueryRow("SELECT role FROM users WHERE username = $1", username).Scan(&role)
	if err != nil {
		return "", err
	}
	return role, nil
}

func Login(c *fiber.Ctx) error {
	// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	var user User
	if err := c.BodyParser(&user); err != nil {
		c.JSON(err)
	}

	// Retrieve the user's hashed password from the database (PostgreSQL)
	// Example: storedPassword := "hashed_password_retrieved_from_database"
	hashedPassword, err := getUserPassword(user.Username)
	if err != nil {
		c.JSON(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		c.JSON("Invalid username or password")
	}

	userRole, err := getUserRole(user.Username)
	if err != nil {
		c.JSON("Invalid username or password")

	}

	fmt.Println(userRole)
	// Include the user's role in the JWT token payload
	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     userRole, // Include the user's role in the token
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Fatal("Error generating token")
	}

	//return c.JSON(fiber.Map{"token": tokenString})
	c.Set("Content-Type", "application/json")
	return c.Status(http.StatusOK).JSON(LoginResponse{
		Status:   http.StatusOK,
		Token:    tokenString,
		Role:     userRole,
		Username: user.Username,
	})

}

func GetUsers(c *fiber.Ctx) error {
	// Extract the JWT token from the request header
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		c.JSON("Authorization token is required")
	}

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conforms to the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return c.JSON("Invalid token")

	}

	// Check if the token is valid
	if !token.Valid {
		return c.JSON("Invalid token")

	}

	// Extract the claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON("Invalid token claims")

	}
	// fmt.Println(token.Claims.(jwt.MapClaims))
	// fmt.Println(claims["role"])
	// Check the user's role from the token
	role, ok := claims["role"].(string)
	if !ok {
		return c.JSON("Invalid user role")

	}

	// Check the user's role and allow access based on the role
	if role == "admin" {
		// Allow admin access
		// Your admin logic goes here

		rows, err := db.Query("SELECT * FROM users;")
		if err != nil {
			log.Fatal(err)
		}

		var users []User

		for rows.Next() {
			var user User
			if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedOn,
				&user.Role); err != nil {
				log.Fatal(err)
			}
			users = append(users, user)
		}

		if rows.Err(); err != nil {
			log.Fatal(err)
		}

		// result, err := json.Marshal(users)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		return c.JSON(fiber.Map{"data": users})

	} else {
		// Allow normal user access
		// Your normal user logic goes here

		return c.JSON("You could not access this section")
	}
}

// func authorize(w http.ResponseWriter, r *http.Request) {
// 	// Extract the JWT token from the request header
// 	w.Header().Set("Content-Type", "application/json")
// 	tokenString := r.Header.Get("Authorization")
// 	fmt.Println("tokenString: ", tokenString)
// 	if tokenString == "" {
// 		http.Error(w, "Authorization token is required", http.StatusUnauthorized)
// 		return
// 	}

// 	// Parse the JWT token
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		// Validate the token with the signing method and key
// 		// ...

// 		return jwtKey, nil
// 	})
// 	if err != nil {
// 		http.Error(w, "Invalid token", http.StatusUnauthorized)
// 		return
// 	}

// 	// Check if the token is valid
// 	if !token.Valid {
// 		http.Error(w, "Invalid token", http.StatusUnauthorized)
// 		return
// 	}

// 	// If the token is valid, respond with a success status
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(tokenString))
// }
