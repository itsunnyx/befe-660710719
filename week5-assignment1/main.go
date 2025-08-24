package main

import (
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
)

type Bookstore struct {
    ID       	string  `json:"id"`
    Name     	string  `json:"name"`
	Category   	string  `json:"category"`
    Price     	int     `json:"price"`
	Author		string	`json:"author"`
}

type Employees struct {
    ID       	string  `json:"id"`
    Name     	string  `json:"name"`
	Position   	string  `json:"position"`
    Salary     	int     `json:"salary"`
}

var bookstores = []Bookstore{
    {
		ID: "1",
		Name: "Clean Code",
		Category: "Programming",
		Price: 650,
		Author: "Robert C. Martin",
	},
	{
		ID: "2",
		Name: "Harry Potter and the Sorcerer's Stone",
		Category: "Fantasy",
		Price: 450,
		Author: "J.K. Rowling",
	},
	{
		ID: "3",
		Name: "Atomic Habits",
		Category: "Self-Improvement",
		Price: 380,
		Author: "James Clear",
	},
}

var employees = []Employees{
	{
		ID:       "B001",
		Name:     "Somsak Chaiyasit",
		Position: "Store Manager",
		Salary:   45000,
	},
	{
		ID:       "B002",
		Name:     "Narumon Prasert",
		Position: "Cashier",
		Salary:   25000,
	},
	{
		ID:       "B003",
		Name:     "Kittipong Wongchai",
		Position: "Sales Assistant",
		Salary:   22000,
	},
	{
		ID:       "B004",
		Name:     "Siriporn Lertdee",
		Position: "Inventory Staff",
		Salary:   24000,
	},
	{
		ID:       "B005",
		Name:     "Anucha Meechai",
		Position: "Book Curator",
		Salary:   30000,
	},
}

func getBookstore(c *gin.Context) {
	bookID := c.Query("id")
	if (bookID != "") {
		filter := []Bookstore{}
		for _, bookstore := range bookstores {
			if (fmt.Sprint(bookstore.ID) == bookID){
				filter = append(filter, bookstore)
			}
		}
		c.JSON(http.StatusOK, filter)
		return
	}
	c.JSON(http.StatusOK, bookstores)
}

func getEmployee(c *gin.Context) {
	employeeID := c.Query("id")
	if (employeeID != "") {
		filter := []Employees{}
		for _, employee := range employees {
			if (fmt.Sprint(employee.ID) == employeeID){
				filter = append(filter, employee)
			}
		}
		c.JSON(http.StatusOK, filter)
		return
	}
	c.JSON(http.StatusOK, employees)
}

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "healthy"})
	})

	api := r.Group("/api/v1")
	{
		api.GET("/bookstores", getBookstore)
		api.GET("/employees", getEmployee)
	}

	r.Run(":8080")
}