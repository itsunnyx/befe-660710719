package main

import (

	_ "week11-assignment/docs"

	"fmt"
	"os"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
    "github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-contrib/cors"

	"strings"
	"strconv"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var db *sql.DB

func initDB() {

    if err := godotenv.Load(); err != nil {
        log.Println("⚠️  .env not found, using system environment variables")
    }

	var err error
	// host := getEnv("DB_HOST", "localhost")
	// name := getEnv("DB_NAME", "bookstore")
	// user := getEnv("DB_USER", "bookstore_user")
	// password := getEnv("DB_PASSWORD", "your_strong_password")
	// port := getEnv("DB_PORT", "5432")

	// conSt := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)
	// fmt.Println(conSt)
    conSt := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err = sql.Open("postgres", conSt)
	if err != nil {
		log.Fatal("failed to open database")
	}

	// กำหนดจำนวน Connection สูงสุด
	db.SetMaxOpenConns(25)

	// กำหนดจำนวน Idle connection สูงสุด	
	db.SetMaxIdleConns(20)

	// กำหนดอายุของ Connection
	db.SetConnMaxLifetime(5 * time.Minute)	

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		log.Fatal("failed to connect database")
	}

	log.Println("succesfully connected to database")
}

type Book struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	ISBN          string    `json:"isbn"`
	Year          int       `json:"year"`
	Price         float64   `json:"price"`

	// ฟิลด์ใหม่
	Category      string    `json:"category"`
	OriginalPrice *float64  `json:"original_price,omitempty"`
	Discount      int       `json:"discount"`
	CoverImage    string    `json:"cover_image"`
	Rating        float64   `json:"rating"`
	ReviewsCount  int       `json:"reviews_count"`
	IsNew         bool      `json:"is_new"`
	Pages         *int      `json:"pages,omitempty"`
	Language      string    `json:"language"`
	Publisher     string    `json:"publisher"`
	Description   string    `json:"description"`

	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type SearchBooksResponse struct {
    Query    string `json:"query"`
    Page     int    `json:"page"`
    PerPage  int    `json:"per_page"`
    Total    int    `json:"total"`
    Results  []Book `json:"results"`
}

// @Summary Get all book
// @Description Get details of book
// @Tags Books
// @Produce  json
// @Success 200  {array}  Book
// @Failure 404  {object}  ErrorResponse
// @Router  /books [get]
func getAllBooks(c *gin.Context) {

    var rows *sql.Rows
    var err error

    // ลูกค้าถาม "มีหนังสืออะไรบ้าง"
    rows, err = db.Query(`SELECT id, title, author, isbn, year, price, 
	category, original_price, discount, cover_image, rating, reviews_count, is_new, pages,
	language, publisher, description, created_at, updated_at FROM books`)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close() // ต้องปิด rows เสมอ เพื่อคืน Connection กลับ pool

    var books []Book
    for rows.Next() {
        var book Book
        err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, 
			&book.Category, &book.OriginalPrice, &book.Discount, &book.CoverImage, &book.Rating, 
			&book.ReviewsCount, &book.IsNew, &book.Pages, &book.Language, &book.Publisher, 
			&book.Description, &book.CreatedAt, &book.UpdatedAt)
        if err != nil {
            // handle error
        }
        books = append(books, book)
    }
	if books == nil {
		books = []Book{}
	}

	yearQuery := c.Query("year")
	if (yearQuery != "") {
		filter := []Book{}
		for _, book := range books {
			if (fmt.Sprint(book.Year) == yearQuery){
				filter = append(filter, book)
			}
		}
		c.JSON(http.StatusOK, filter)
		return
	}

	categoryQuery := c.Query("category")
	if (categoryQuery != "") {
		filter := []Book{}
		for _, book := range books {
			if (fmt.Sprint(book.Category) == categoryQuery){
				filter = append(filter, book)
			}
		}
		c.JSON(http.StatusOK, filter)
		return
	}

	c.JSON(http.StatusOK, books)
}

// @Summary Get book by ID
// @Description Get details of a book by ID
// @Tags Books
// @Produce  json
// @Param   id   path      int     true  "Book ID"
// @Success 200  {object}  Book
// @Failure 404  {object}  ErrorResponse
// @Router  /books/{id} [get]  
func getBook(c *gin.Context) {
    id := c.Param("id")
    var book Book

    // QueryRow ใช้เมื่อคาดว่าจะได้ผลลัพธ์ 0 หรือ 1 แถว
    err := db.QueryRow("SELECT id, title, author FROM books WHERE id = $1", id).
        Scan(&book.ID, &book.Title, &book.Author)

    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
        return
    } else if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, book)
}

func getNewBooks(c *gin.Context) {

    var rows *sql.Rows
    var err error

    // ลูกค้าถาม "มีหนังสืออะไรบ้าง"
    rows, err = db.Query(`SELECT id, title, author, isbn, year, price, 
	category, original_price, discount, cover_image, rating, reviews_count, is_new, pages,
	language, publisher, description, created_at, updated_at FROM books
	where is_new=true`)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close() // ต้องปิด rows เสมอ เพื่อคืน Connection กลับ pool

    var books []Book
    for rows.Next() {
        var book Book
        err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, 
			&book.Category, &book.OriginalPrice, &book.Discount, &book.CoverImage, &book.Rating, 
			&book.ReviewsCount, &book.IsNew, &book.Pages, &book.Language, &book.Publisher, 
			&book.Description, &book.CreatedAt, &book.UpdatedAt)
        if err != nil {
            // handle error
        }
        books = append(books, book)
    }
	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)
}

func getBookCategories(c *gin.Context) {
	var book Book
	c.JSON(http.StatusOK, book)
}

func searchBooks(c *gin.Context) {
    q := strings.TrimSpace(c.Query("q"))
    if q == "" {
        c.JSON(http.StatusBadRequest, ErrorResponse{Message: "q is required"})
        return
    }

    page := toIntDefault(c.Query("page"), 1)
    if page < 1 { page = 1 }
    per := toIntDefault(c.Query("per_page"), 20)
    if per < 1 { per = 20 }
    if per > 100 { per = 100 }
    offset := (page - 1) * per

    // นับ total
    var total int
    err := db.QueryRow(`
        SELECT COUNT(*) 
        FROM books
        WHERE title ILIKE '%' || $1 || '%'
           OR author ILIKE '%' || $1 || '%'
           OR description ILIKE '%' || $1 || '%'
           OR isbn ILIKE '%' || $1 || '%'
    `, q).Scan(&total)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
        return
    }

    // ถ้าอยากบูสต์ความเกี่ยวข้องแบบเร็ว ๆ ใช้ pg_trgm similarity(title, $1)
    // แนะนำเปิด EXTENSION และทำดัชนี (ดูด้านล่าง)
    rows, err := db.Query(`
		SELECT
			id, title, author, isbn, year, price,
			category, original_price, discount, cover_image,
			rating, reviews_count, is_new, pages,
			language, publisher, description,
			created_at, updated_at
		FROM books
		WHERE title ILIKE '%' || $1 || '%'
		OR author ILIKE '%' || $1 || '%'
		OR description ILIKE '%' || $1 || '%'
		OR isbn ILIKE '%' || $1 || '%'
		ORDER BY rating DESC, reviews_count DESC
		LIMIT $2 OFFSET $3
	`, q, per, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
        return
    }
    defer rows.Close()

    results := make([]Book, 0, per)
    for rows.Next() {
    var b Book
    err := rows.Scan(
        &b.ID, &b.Title, &b.Author, &b.ISBN, &b.Year, &b.Price,
        &b.Category, &b.OriginalPrice, &b.Discount, &b.CoverImage,
        &b.Rating, &b.ReviewsCount, &b.IsNew, &b.Pages,
        &b.Language, &b.Publisher, &b.Description,
        &b.CreatedAt, &b.UpdatedAt,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    results = append(results, b)
}
    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
        return
    }

    c.JSON(http.StatusOK, SearchBooksResponse{
        Query:   q,
        Page:    page,
        PerPage: per,
        Total:   total,
        Results: results,
    })
}

func toIntDefault(s string, def int) int {
    if s == "" { return def }
    v, err := strconv.Atoi(s)
    if err != nil { return def }
    return v
}

func getDiscountedBook(c *gin.Context){
	var rows *sql.Rows
    var err error

    // ลูกค้าถาม "มีหนังสืออะไรบ้าง"
    rows, err = db.Query(`SELECT id, title, author, isbn, year, price, 
	category, original_price, discount, cover_image, rating, reviews_count, is_new, pages,
	language, publisher, description, created_at, updated_at FROM books
	where discount>0`)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close() // ต้องปิด rows เสมอ เพื่อคืน Connection กลับ pool

    var books []Book
    for rows.Next() {
        var book Book
        err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, 
			&book.Category, &book.OriginalPrice, &book.Discount, &book.CoverImage, &book.Rating, 
			&book.ReviewsCount, &book.IsNew, &book.Pages, &book.Language, &book.Publisher, 
			&book.Description, &book.CreatedAt, &book.UpdatedAt)
        if err != nil {
            // handle error
        }
        books = append(books, book)
    }
	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)
}

func getFeaturedBook(c *gin.Context){
	var rows *sql.Rows
    var err error

    // ลูกค้าถาม "มีหนังสืออะไรบ้าง"
    rows, err = db.Query(`SELECT id, title, author, isbn, year, price, 
	category, original_price, discount, cover_image, rating, reviews_count, is_new, pages,
	language, publisher, description, created_at, updated_at FROM books
	where rating>4.5`)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close() // ต้องปิด rows เสมอ เพื่อคืน Connection กลับ pool

    var books []Book
    for rows.Next() {
        var book Book
        err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, 
			&book.Category, &book.OriginalPrice, &book.Discount, &book.CoverImage, &book.Rating, 
			&book.ReviewsCount, &book.IsNew, &book.Pages, &book.Language, &book.Publisher, 
			&book.Description, &book.CreatedAt, &book.UpdatedAt)
        if err != nil {
            // handle error
        }
        books = append(books, book)
    }
	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)
}

// @Summary Create a book by ID
// @Description Create book details (title, author, isbn, year, price) by book ID
// @Tags Books
// @Produce  json
// @Param   book  body      Book    true   "Create book data"
// @Success 200  {object}  Book
// @Failure 404  {object}  ErrorResponse
// @Router  /books [post]  
func createBook(c *gin.Context) {
    var newBook Book

    if err := c.ShouldBindJSON(&newBook); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // ใช้ RETURNING เพื่อดึงค่าที่ database generate (id, timestamps)
    var id int
    var created_At, updated_At time.Time

    err := db.QueryRow(
        `INSERT INTO books (title, author, isbn, year, price)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, created_at, updated_at`,
        newBook.Title, newBook.Author, newBook.ISBN, newBook.Year, newBook.Price,
    ).Scan(&id, &created_At, &updated_At)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    newBook.ID = id
    newBook.CreatedAt = created_At
    newBook.UpdatedAt = updated_At

    c.JSON(http.StatusCreated, newBook) // ใช้ 201 Created
}

// @Summary Update a book by ID
// @Description Update book details (title, author, isbn, year, price) by book ID
// @Tags Books
// @Produce  json
// @Param   id   path      int     true  "Book ID"
// @Param   book  body      Book    true   "Updated book data"
// @Success 200  {object}  Book
// @Failure 404  {object}  ErrorResponse
// @Router  /books/{id} [put]  
func updateBook(c *gin.Context) {
    id := c.Param("id")
    var updateBook Book

    if err := c.ShouldBindJSON(&updateBook); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var updatedAt time.Time
	var return_id int
    err := db.QueryRow(
        `UPDATE books
         SET title = $1, author = $2, isbn = $3, year = $4, price = $5
         WHERE id = $6
         RETURNING id, updated_at`,
        updateBook.Title, updateBook.Author, updateBook.ISBN,
        updateBook.Year, updateBook.Price, id,
    ).Scan(&return_id, &updatedAt)

    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
        return
    } else if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
	updateBook.UpdatedAt = updatedAt
	updateBook.ID = return_id
	c.JSON(http.StatusOK, updateBook)
}

// @Summary Delete a book by ID
// @Description Delete book details by book ID
// @Tags Books
// @Produce  json
// @Param   id   path      int     true  "Book ID"
// @Success 200  {object}  Book
// @Failure 404  {object}  ErrorResponse
// @Router  /books/{id} [delete]  
func deleteBook(c *gin.Context) {
    id := c.Param("id")

    result, err := db.Exec("DELETE FROM books WHERE id = $1", id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if rowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
}

// @title           Simple API Example
// @version         1.0
// @description     This is a simple example of using Gin with Swagger.
// @host            localhost:8080
// @BasePath        /api/v1
func main(){
	initDB()
	defer db.Close()
	
	r := gin.Default()
	r.Use(cors.Default())


	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context){
		err := db.Ping()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "unhealthy", "error":err})
			return
		}
		c.JSON(200, gin.H{"message" : "healthy"})
	})

	api := r.Group("/api/v1")
	{
		api.GET("/books", getAllBooks)
		api.GET("/books/:id", getBook)

		api.GET("/books/new", getNewBooks)
		api.GET("/categories", getBookCategories)
		api.GET("/books/search", searchBooks)
		api.GET("/books/featured", getFeaturedBook)
		api.GET("/books/discounted", getDiscountedBook)

		api.POST("/books", createBook)
		api.PUT("/books/:id", updateBook)

		api.DELETE("/books/:id", deleteBook)
	}

	r.Run(":8080")
}