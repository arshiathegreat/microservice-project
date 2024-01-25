package storeService


import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cs_student_uni/microservice_store/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type Product struct {
	ID          int32   `json:"product_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	MadeIn      string  `json:"made_in"`
}

var db *sql.DB
var jwtKey = []byte("eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiQWRtaW4iLCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6InN0dWRlbnRfY3MiLCJleHAiOjE3MzczMTcyMjksImlhdCI6MTcwNTY5NDgyOX0.U9ahjosw9nIhDY2c3NBmgjf4hB89RUzLG_ROIrtIEnM")

func Init() {
	db = database.GetDB()
	// http.HandleFunc("/products", getProducts)
	// http.HandleFunc("/createProduct", createProduct)
	// http.HandleFunc("/sellProduct", sellProduct)
	//http.HandleFunc("/sellProduct", authorize)
}

func GetProducts(c *fiber.Ctx) error {

	//err := json.NewDecoder(r.Body).Decode(&products)
	//if err != nil {
	//	panic(err)
	//}

	rows, err := db.Query("SELECT * FROM products;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// An product slice to hold data from returned rows.
	var products []Product

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Title, &product.Description, &product.MadeIn,
			&product.Price); err != nil {
			log.Fatal("Error next", err)
		}
		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		log.Fatal("Error rows:", err)
	}

	// result, err := json.Marshal(products)
	// if err != nil {
	// 	panic(err)
	// }

	return c.JSON(fiber.Map{"data": products})
}

func saveProductToDB(title string, description string, price float64, made_in string) error {
	_, err := db.Exec("INSERT INTO products (title, description, price, made_in) VALUES ($1, $2, $3, $4)", title, description, price, made_in)
	if err != nil {
		fmt.Println("Error message: ", err)
		return err
	}
	return nil
}

func CreateProduct(c *fiber.Ctx) error {
	var product Product
	if err := c.BodyParser(&product); err != nil {
		return c.JSON(err)
	}

	err := saveProductToDB(product.Title, product.Description, product.Price, product.MadeIn)
	if err != nil {
		return c.JSON("Error saving product to database")
	}

	return c.JSON(fiber.Map{"data": product})
}

func SellProduct(c *fiber.Ctx) error {

	// Extract the JWT token from the request header
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.JSON("Authorization token is required")
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
	fmt.Println(token.Claims.(jwt.MapClaims))
	fmt.Println(claims["role"])
	// Check the user's role from the token
	role, ok := claims["role"].(string)
	if !ok {
		return c.JSON("Invalid user role")
	}

	var products []Product
	err = c.BodyParser(&products)
	if err != nil {
		return c.JSON(err)
	}

	result, err := json.Marshal(products)
	if err != nil {
		panic(err)
	}

	// Check the user's role and allow access based on the role
	if role == "customer" {
		var exists bool
		// "SELECT EXISTS(SELECT title FROM products WHERE title = ?);"
		for _, v := range products {
			err = db.QueryRow(`SELECT
			EXISTS(
				SELECT
					title
				FROM
					products
				WHERE
					title = $1
			);`, v.Title).Scan(&exists)
			if err != nil {
				panic(err.Error())
			}

		}
		if exists {
			return c.JSON(fiber.Map{"products": result, "message": "Sell logic executed successfully"})
		} else {

			return c.JSON(fiber.Map{"message": "The product does not exist in the database"})
		}

	} else {
		// Allow normal user access
		// Your normal user logic goes here
		return c.JSON(fiber.Map{"message": "Sell logic executed failed"})

	}

}
