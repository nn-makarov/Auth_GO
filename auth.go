package main

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
    "os"
)

var jwtSecret = []byte("my-super-secret-key-12345")

type Claims struct {
    UserID uint `json:"user_id"`
    jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func GenerateToken(userID uint) (string, error) {
    claims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (uint, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    
    if err != nil {
        return 0, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims.UserID, nil
    }
    
    return 0, errors.New("invalid token")
}

func CheckPermission(db *gorm.DB, user *User, resourceName string, action string) bool {
    if user.Role == "admin" {
        return true
    }
    
    var resource Resource
    if err := db.Where("name = ?", resourceName).First(&resource).Error; err != nil {
        return false
    }
    
    var permission Permission
    err := db.Where("role = ? AND resource_id = ?", user.Role, resource.ID).First(&permission).Error
    if err != nil {
        return false
    }
    
    switch action {
    case "view":
        return permission.CanView
    case "create":
        return permission.CanCreate
    case "edit":
        return permission.CanEdit
    case "delete":
        return permission.CanDelete
    default:
        return false
    }
}

func InitDB() *gorm.DB {
    // Пробуем открыть или создать БД
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        // Если не открывается, пробуем удалить и создать заново
        os.Remove("test.db")
        db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
        if err != nil {
            panic("failed to connect database: " + err.Error())
        }
    }
    
    // Автомиграция
    if err := db.AutoMigrate(&User{}, &Resource{}, &Permission{}); err != nil {
        panic("failed to migrate database: " + err.Error())
    }
    
    // Создаем ресурсы
    resources := []string{"products", "orders", "users"}
    for _, name := range resources {
        var resource Resource
        result := db.FirstOrCreate(&resource, Resource{Name: name})
        if result.Error != nil {
            panic("failed to create resource: " + result.Error.Error())
        }
    }
    
    // Создаем админа, если нет
    var admin User
    if err := db.Where("email = ?", "admin@test.com").First(&admin).Error; err != nil {
        hashedPassword, err := HashPassword("admin123")
        if err != nil {
            panic("failed to hash password: " + err.Error())
        }
        admin = User{
            Email:        "admin@test.com",
            PasswordHash: hashedPassword,
            FirstName:    "Admin",
            LastName:     "User",
            Role:         "admin",
            IsActive:     true,
        }
        if err := db.Create(&admin).Error; err != nil {
            panic("failed to create admin: " + err.Error())
        }
    }
    
    return db
}