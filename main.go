package main


import (
    "database/sql"
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    _ "github.com/lib/pq"
)

var db *sql.DB
var jwtKey = []byte("secret_key")

type User struct {
    ID       int
    Name     string
    Email    string
    Password string
    Role     string
}

type Car struct {
    ID        int
    UserID    int
    Marca     string
    Modelo    string
    Ano       int
    Valor     float64
}


// =============================
// CONEXAO DB
// =============================

func initDB() {
    connStr := "host=db user=postgres password=postgres dbname=goapi port=5432 sslmode=disable"
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
}

// =============================
// JWT
// =============================

func generateToken(user User) (string, error) {
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "role":    user.Role,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

// =============================
// MIDDLEWARE
// =============================

func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenStr := c.GetHeader("Authorization")

        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        claims := token.Claims.(jwt.MapClaims)
        c.Set("user_id", int(claims["user_id"].(float64)))
        c.Set("role", claims["role"].(string))

        c.Next()
    }
}

func adminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, _ := c.Get("role")
        if role != "admin" {
            c.AbortWithStatus(http.StatusForbidden)
            return
        }
        c.Next()
    }
}

// =============================
// HANDLERS
// =============================

func register(c *gin.Context) {
    var u User
    c.BindJSON(&u)

    err := db.QueryRow(
        "INSERT INTO users(name,email,password,role) VALUES($1,$2,$3,$4) RETURNING id",
        u.Name, u.Email, u.Password, u.Role,
    ).Scan(&u.ID)

    if err != nil {
        c.JSON(500, err)
        return
    }

    c.JSON(200, u)
}

func login(c *gin.Context) {
    var input User
    c.BindJSON(&input)

    var user User
    err := db.QueryRow(
        "SELECT id, password, role FROM users WHERE email=$1",
        input.Email,
    ).Scan(&user.ID, &user.Password, &user.Role)

    if err != nil || user.Password != input.Password {
        c.JSON(401, gin.H{"error": "invalid credentials"})
        return
    }

    token, _ := generateToken(user)
    c.JSON(200, gin.H{"token": token})
}

func createCar(c *gin.Context) {
    var car Car
    c.BindJSON(&car)

    userID, _ := c.Get("user_id")

    err := db.QueryRow(
        "INSERT INTO cars(user_id,marca,modelo,ano,valor) VALUES($1,$2,$3,$4,$5) RETURNING id",
        userID, car.Marca, car.Modelo, car.Ano, car.Valor,
    ).Scan(&car.ID)

    if err != nil {
        c.JSON(500, err)
        return
    }

    c.JSON(200, car)
}

// =============================
// MAIN
// =============================

func main() {
    initDB()

    r := gin.Default()

    r.POST("/register", register)
    r.POST("/login", login)

    auth := r.Group("/")
    auth.Use(authMiddleware())

    auth.POST("/cars", createCar)

    admin := auth.Group("/admin")
    admin.Use(adminMiddleware())
    admin.POST("/users", register)

    r.Run(":8080")
}
package main

import (

func register(c *gin.Context) {
    var u User
    c.BindJSON(&u)

    err := db.QueryRow(
        "INSERT INTO users(name,email,password,role) VALUES($1,$2,$3,$4) RETURNING id",
        u.Name, u.Email, u.Password, u.Role,
    ).Scan(&u.ID)

    if err != nil {
        c.JSON(500, err)
        return
    }

    c.JSON(200, u)
}

func login(c *gin.Context) {
    var input User
    c.BindJSON(&input)

    var user User
    err := db.QueryRow(
        "SELECT id, password, role FROM users WHERE email=$1",
        input.Email,
    ).Scan(&user.ID, &user.Password, &user.Role)

    if err != nil || user.Password != input.Password {
        c.JSON(401, gin.H{"error": "invalid credentials"})
        return
    }

    token, _ := generateToken(user)
    c.JSON(200, gin.H{"token": token})
}

func createCar(c *gin.Context) {
    var car Car
    c.BindJSON(&car)

    userID, _ := c.Get("user_id")

    err := db.QueryRow(
        "INSERT INTO cars(user_id,marca,modelo,ano,valor) VALUES($1,$2,$3,$4,$5) RETURNING id",
        userID, car.Marca, car.Modelo, car.Ano, car.Valor,
    ).Scan(&car.ID)

    if err != nil {
        c.JSON(500, err)
        return
    }

    c.JSON(200, car)
}

// =============================
// MAIN
// =============================

func main() {
    initDB()

    r := gin.Default()

    r.POST("/register", register)
    r.POST("/login", login)

    auth := r.Group("/")
    auth.Use(authMiddleware())

    auth.POST("/cars", createCar)

    admin := auth.Group("/admin")
    admin.Use(adminMiddleware())
    admin.POST("/users", register)

    r.Run(":8080")
}
package main

import (

func register(c *gin.Context) {
    var u User
    c.BindJSON(&u)

    err := db.QueryRow(
        "INSERT INTO users(name,email,password,role) VALUES($1,$2,$3,$4) RETURNING id",
        u.Name, u.Email, u.Password, u.Role,
    ).Scan(&u.ID)

    if err != nil {
        c.JSON(500, err)
        return
    }

    c.JSON(200, u)
}

func login(c *gin.Context) {
    var input User
    c.BindJSON(&input)

    var user User
    err := db.QueryRow(
        "SELECT id, password, role FROM users WHERE email=$1",
        input.Email,
    ).Scan(&user.ID, &user.Password, &user.Role)

    if err != nil || user.Password != input.Password {
        c.JSON(401, gin.H{"error": "invalid credentials"})
        return
    }

    token, _ := generateToken(user)
    c.JSON(200, gin.H{"token": token})
}

func createCar(c *gin.Context) {
    var car Car
    c.BindJSON(&car)

    userID, _ := c.Get("user_id")

    err := db.QueryRow(
        "INSERT INTO cars(user_id,marca,modelo,ano,valor) VALUES($1,$2,$3,$4,$5) RETURNING id",
        userID, car.Marca, car.Modelo, car.Ano, car.Valor,
    ).Scan(&car.ID)

    if err != nil {
        c.JSON(500, err)
        return
    }

    c.JSON(200, car)
}

// =============================
// MAIN
// =============================

func main() {
    initDB()

    r := gin.Default()

    r.POST("/register", register)
    r.POST("/login", login)

    auth := r.Group("/")
    auth.Use(authMiddleware())

    auth.POST("/cars", createCar)

    admin := auth.Group("/admin")
    admin.Use(adminMiddleware())
    admin.POST("/users", register)

    r.Run(":8080")
}
package main

import (

func register(c *gin.Context) {
    var u User
    c.BindJSON(&u)

    err := db.QueryRow(
        "INSERT INTO users(name,email,password,role) VALUES($1,$2,$3,$4) RETURNING id",
        u.Name, u.Email, u.Password, u.Role,
    ).Scan(&u.ID)

    if err != nil {
        c.JSON(500, err)
        return
    }

    c.JSON(200, u)
}

func login(c *gin.Context) {
    var input User
    c.BindJSON(&input)

    var user User
    err := db.QueryRow(
        "SELECT id, password, role FROM users WHERE email=$1",
        input.Email,
    ).Scan(&user.ID, &user.Password, &user.Role)

    if err != nil || user.Password != input.Password {
        c.JSON(401, gin.H{"error": "invalid credentials"})
        return
    }

    token, _ := generateToken(user)
    c.JSON(200, gin.H{"token": token})
}

func createCar(c *gin.Context) {
    var car Car
    c.BindJSON(&car)

    userID, _ := c.Get("user_id")

    err := db.QueryRow(
        "INSERT INTO cars(user_id,marca,modelo,ano,valor) VALUES($1,$2,$3,$4,$5) RETURNING id",
        userID, car.Marca, car.Modelo, car.Ano, car.Valor,
    ).Scan(&car.ID)

    if err != nil {
        c.JSON(500, err)
        return
    }

    c.JSON(200, car)
}

// =============================
// MAIN
// =============================

func main() {
    initDB()

    r := gin.Default()

    r.POST("/register", register)
    r.POST("/login", login)

    auth := r.Group("/")
    auth.Use(authMiddleware())

    auth.POST("/cars", createCar)

    admin := auth.Group("/admin")
    admin.Use(adminMiddleware())
    admin.POST("/users", register)

    r.Run(":8080")
}
package main

import (

func register(c *gin.Context) {
    var u User
    c.BindJSON(&u)

    err := db.QueryRow(
        "INSERT INTO users(name,email,password,role) VALUES($1,$2,$3,$4) RETURNING id",
        u.Name, u.Email, u.Password, u.Role,
    ).Scan(&u.ID)

    if err != nil {
        c.JSON(500, err)
        return
    }

    c.JSON(200, u)
}

func login(c *gin.Context) {
    var input User
    c.BindJSON(&input)

    var user User
    err := db.QueryRow(
        "SELECT id, password, role FROM users WHERE email=$1",
        input.Email,
    ).Scan(&user.ID, &user.Password, &user.Role)

    if err != nil || user.Password != input.Password {
        c.JSON(401, gin.H{"error": "invalid credentials"})
        return
    }

    token, _ := generateToken(user)
    c.JSON(200, gin.H{"token": token})
}

func createCar(c *gin.Context) {
    var car Car
    c.BindJSON(&car)

    userID, _ := c.Get("user_id")

    err := db.QueryRow(
        "INSERT INTO cars(user_id,marca,modelo,ano,valor) VALUES($1,$2,$3,$4,$5) RETURNING id",
        userID, car.Marca, car.Modelo, car.Ano, car.Valor,
    ).Scan(&car.ID)

    if err != nil {
        c.JSON(500, err)
        return
    }

    c.JSON(200, car)
}

// =============================
// MAIN
// =============================

func main() {
    initDB()

    r := gin.Default()

    r.POST("/register", register)
    r.POST("/login", login)

    auth := r.Group("/")
    auth.Use(authMiddleware())

    auth.POST("/cars", createCar)

    admin := auth.Group("/admin")
    admin.Use(adminMiddleware())
    admin.POST("/users", register)

    r.Run(":8080")
}
package main

import (

func register(c *gin.Context) {
    var u User
    c.BindJSON(&u)

    err := db.QueryRow(
        "INSERT INTO users(name,email,password,role) VALUES($1,$2,$3,$4) RETURNING id",
        u.Name, u.Email, u.Password, u.Role,
    ).Scan(&u.ID)

    if err != nil {
        c.JSON(500, err)
        return
    }

    c.JSON(200, u)
}

func login(c *gin.Context) {
    var input User
    c.BindJSON(&input)

    var user User
    err := db.QueryRow(
        "SELECT id, password, role FROM users WHERE email=$1",
        input.Email,
    ).Scan(&user.ID, &user.Password, &user.Role)

    if err != nil || user.Password != input.Password {
        c.JSON(401, gin.H{"error": "invalid credentials"})
        return
    }

    token, _ := generateToken(user)
    c.JSON(200, gin.H{"token": token})
}

func createCar(c *gin.Context) {
    var car Car
    c.BindJSON(&car)

    userID, _ := c.Get("user_id")

    err := db.QueryRow(
        "INSERT INTO cars(user_id,marca,modelo,ano,valor) VALUES($1,$2,$3,$4,$5) RETURNING id",
        userID, car.Marca, car.Modelo, car.Ano, car.Valor,
    ).Scan(&car.ID)

    if err != nil {
        c.JSON(500, err)
        return
    }

    c.JSON(200, car)
}

// =============================
// MAIN
// =============================

func main() {
    initDB()

    r := gin.Default()

    r.POST("/register", register)
    r.POST("/login", login)

    auth := r.Group("/")
    auth.Use(authMiddleware())

    auth.POST("/cars", createCar)

    admin := auth.Group("/admin")
    admin.Use(adminMiddleware())
    admin.POST("/users", register)

    r.Run(":8080")
}
package main

import (

func register(c *gin.Context) {
    var u User
    c.BindJSON(&u)

    err := db.QueryRow(
        "INSERT INTO users(name,email,password,role) VALUES($1,$2,$3,$4) RETURNING id",
        u.Name, u.Email, u.Password, u.Role,
    ).Scan(&u.ID)

    if err != nil {
        c.JSON(500, err)
        return
    }

    c.JSON(200, u)
}

func login(c *gin.Context) {
    var input User
    c.BindJSON(&input)

    var user User
    err := db.QueryRow(
        "SELECT id, password, role FROM users WHERE email=$1",
        input.Email,
    ).Scan(&user.ID, &user.Password, &user.Role)

    if err != nil || user.Password != input.Password {
        c.JSON(401, gin.H{"error": "invalid credentials"})
        return
    }

    token, _ := generateToken(user)
    c.JSON(200, gin.H{"token": token})
}

func createCar(c *gin.Context) {
    var car Car
    c.BindJSON(&car)

    userID, _ := c.Get("user_id")

    err := db.QueryRow(
        "INSERT INTO cars(user_id,marca,modelo,ano,valor) VALUES($1,$2,$3,$4,$5) RETURNING id",
        userID, car.Marca, car.Modelo, car.Ano, car.Valor,
    ).Scan(&car.ID)

    if err != nil {
        c.JSON(500, err)
        return
    }

    c.JSON(200, car)
}

// =============================
// MAIN
// =============================

func main() {
    initDB()

    r := gin.Default()

    r.POST("/register", register)
    r.POST("/login", login)

    auth := r.Group("/")
    auth.Use(authMiddleware())

    auth.POST("/cars", createCar)

    admin := auth.Group("/admin")
    admin.Use(adminMiddleware())
    admin.POST("/users", register)

    r.Run(":8080")
}
package main

import (

func register(c *gin.Context) {
    var u User
    c.BindJSON(&u)

    err := db.QueryRow(
        "INSERT INTO users(name,email,password,role) VALUES($1,$2,$3,$4) RETURNING id",
        u.Name, u.Email, u.Password, u.Role,
    ).Scan(&u.ID)

    if err != nil {
        c.JSON(500, err)
        return
    }

    c.JSON(200, u)
}

func login(c *gin.Context) {
    var input User
    c.BindJSON(&input)

    var user User
    err := db.QueryRow(
        "SELECT id, password, role FROM users WHERE email=$1",
        input.Email,
    ).Scan(&user.ID, &user.Password, &user.Role)

    if err != nil || user.Password != input.Password {
        c.JSON(401, gin.H{"error": "invalid credentials"})
        return
    }

    token, _ := generateToken(user)
    c.JSON(200, gin.H{"token": token})
}

func createCar(c *gin.Context) {
    var car Car
    c.BindJSON(&car)

    userID, _ := c.Get("user_id")

    err := db.QueryRow(
        "INSERT INTO cars(user_id,marca,modelo,ano,valor) VALUES($1,$2,$3,$4,$5) RETURNING id",
        userID, car.Marca, car.Modelo, car.Ano, car.Valor,
    ).Scan(&car.ID)

    if err != nil {
        c.JSON(500, err)
        return
    }

    c.JSON(200, car)
}

// =============================
// MAIN
// =============================

func main() {
    initDB()

    r := gin.Default()

    r.POST("/register", register)
    r.POST("/login", login)

    auth := r.Group("/")
    auth.Use(authMiddleware())

    auth.POST("/cars", createCar)

    admin := auth.Group("/admin")
    admin.Use(adminMiddleware())
    admin.POST("/users", register)

    r.Run(":8080")
}

