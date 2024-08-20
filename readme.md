# WarungJWT PostgreSQL Migration Guide

This guide provides step-by-step instructions for migrating your project from MySQL to PostgreSQL using GORM.

## Prerequisites

Before you start, ensure that you have the following installed:
- Go 1.22.1 or later
- PostgreSQL database
- GORM library

## 1. Install PostgreSQL Driver

Replace the MySQL driver with the PostgreSQL driver by running the following command:

```bash
go get gorm.io/driver/postgres
go mod tidy
```

## 2. Update Database Configuration

In your config/database.go file, update the database connection to use PostgreSQL:

```go
package config

import (
    "fmt"
    "log"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/joho/godotenv"
)

var DB *gorm.DB

func InitDatabase() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbHost, dbUser, dbPassword, dbName, dbPort)
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to the database:", err)
    }

    log.Println("Successfully connected to the database")
}

```

## 3. Update Models to Use GORM

Update your models to include GORM tags and adjust the struct fields to work with PostgreSQL.

A. models/product.go

```go
package models

import "gorm.io/gorm"

type Product struct {
    gorm.Model
    Name        string `json:"name"`
    Code        string `json:"code"`
    Stock       int    `json:"stock"`
    Description string `json:"description"`
    Status      string `json:"status"`
}

```

B. models/user.go

```go
package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Email    string `json:"email" gorm:"unique"`
    Password string `json:"password"`
    FullName string `json:"full_name"`
    Role     string `json:"role"`
}

```

## 4. Refactor Handlers to use GORM

Refactor your handlers to use GORM for database operations.

A. handlers/user.go

```go
package handlers

import (
    "net/http"
    "warungjwt_postgre/config"
    "warungjwt_postgre/models"

    "github.com/labstack/echo/v4"
    "golang.org/x/crypto/bcrypt"
)

func Register(c echo.Context) error {
    email := c.FormValue("email")
    password := c.FormValue("password")
    fullName := c.FormValue("full_name")
    role := c.FormValue("role")

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

    user := models.User{
        Email:    email,
        Password: string(hashedPassword),
        FullName: fullName,
        Role:     role,
    }

    if err := config.DB.Create(&user).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to register user"})
    }

    token, err := generateJWT(user.Email, user.Role)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
    }

    return c.JSON(http.StatusCreated, map[string]string{
        "message": "User successfully registered",
        "token":   token,
    })
}

```

B. handlers/product.go

```bash
package handlers

import (
    "net/http"
    "strconv"
    "warungjwt_postgre/config"
    "warungjwt_postgre/models"

    "github.com/labstack/echo/v4"
)

func GetProducts(c echo.Context) error {
    var products []models.Product
    if err := config.DB.Find(&products).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }
    return c.JSON(http.StatusOK, products)
}

func GetProduct(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
    }

    var product models.Product
    if err := config.DB.First(&product, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
    }

    return c.JSON(http.StatusOK, product)
}

func CreateProduct(c echo.Context) error {
    var product models.Product
    if err := c.Bind(&product); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
    }

    if err := config.DB.Create(&product).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create product"})
    }

    return c.JSON(http.StatusCreated, product)
}

func UpdateProduct(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
    }

    var product models.Product
    if err := config.DB.First(&product, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
    }

    if err := c.Bind(&product); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
    }

    if err := config.DB.Save(&product).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update product"})
    }

    return c.JSON(http.StatusOK, product)
}

func DeleteProduct(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
    }

    if err := config.DB.Delete(&models.Product{}, id).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete product"})
    }

    return c.NoContent(http.StatusNoContent)
}

```

## 5. Update .env 

Update your .env file to reflect PostgreSQL configuration:

```bash
DB_USER=your_postgres_user
DB_PASSWORD=your_postgres_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=your_postgres_dbname
```




## 6. Migrate your db

If necessary, migrate your MySQL database to PostgreSQL using a migration tool like pgloader or perform the migration manually.

1. Buat nama db di Postgres 

```sql
CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  code VARCHAR(100) NOT NULL,
  stock INT NOT NULL,
  description TEXT,
  status VARCHAR(20) NOT NULL
);

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  full_name VARCHAR(50) NOT NULL,
  role VARCHAR(20) DEFAULT 'staff',
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO products (name, code, stock, description, status) VALUES
('Product 1', 'P001', 10, 'Description for Product 1', 'active'),
('Product 2', 'P002', 5, 'Description for Product 2', 'broken'),
('Product 3', 'P003', 15, 'Description for Product 3', 'active'),
('Product 99', 'P099', 10, 'Description for Product 99', 'active'),
('Product 98', 'P098', 10, 'Description for Product 98', 'active'),
('Product 77', 'P077', 10, 'Description for Product 77', 'active');
```

atau jika menggunakan pgloader : 

Lakukan Migrasi Menggunakan pgloader 

contoh : 
```bash
pgloader mysql://your_mysql_user:your_mysql_password@localhost/simple_webserver postgresql://your_postgres_user:your_postgres_password@localhost/warung
```



## 7. Run the App

After completing the above steps, you can run your application:

```bash
go run main.go
```