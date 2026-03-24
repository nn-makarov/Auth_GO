package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req RegisterRequest
        
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        var existing User
        if err := db.Where("email = ?", req.Email).First(&existing).Error; err == nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
            return
        }
        
        hashedPassword, err := HashPassword(req.Password)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
            return
        }
        
        user := User{
            Email:        req.Email,
            PasswordHash: hashedPassword,
            FirstName:    req.FirstName,
            LastName:     req.LastName,
            Role:         "user",
            IsActive:     true,
        }
        
        if err := db.Create(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
            return
        }
        
        c.JSON(http.StatusCreated, gin.H{"message": "user created", "email": user.Email})
    }
}

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req LoginRequest
        
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        var user User
        if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
            return
        }
        
        if !CheckPasswordHash(req.Password, user.PasswordHash) {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
            return
        }
        
        if !user.IsActive {
            c.JSON(http.StatusForbidden, gin.H{"error": "account is disabled"})
            return
        }
        
        token, err := GenerateToken(user.ID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
            return
        }
        
        c.JSON(http.StatusOK, LoginResponse{
            Token: token,
            Email: user.Email,
            Role:  user.Role,
        })
    }
}

func ProfileHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("user_id")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            return
        }
        
        var user User
        if err := db.First(&user, userID).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "id":         user.ID,
            "email":      user.Email,
            "first_name": user.FirstName,
            "last_name":  user.LastName,
            "role":       user.Role,
        })
    }
}

func ProductsHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := c.MustGet("user").(User)
        
        if !CheckPermission(db, &user, "products", "view") {
            c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "message":  "list of products",
            "products": []string{"product1", "product2"},
        })
    }
}

func CreatePermissionHandler(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := c.MustGet("user").(User)
        if user.Role != "admin" {
            c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
            return
        }
        
        var req CreatePermissionRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        
        var resource Resource
        if err := db.Where("name = ?", req.Resource).First(&resource).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
            return
        }
        
        permission := Permission{
            Role:       req.Role,
            ResourceID: resource.ID,
            CanView:    req.CanView,
            CanCreate:  req.CanCreate,
            CanEdit:    req.CanEdit,
            CanDelete:  req.CanDelete,
        }
        
        db.Where(Permission{Role: req.Role, ResourceID: resource.ID}).
           Assign(permission).
           FirstOrCreate(&permission)
        
        c.JSON(http.StatusOK, gin.H{"message": "permission created", "permission": permission})
    }
}

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
            c.Abort()
            return
        }
        
        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }
        
        userID, err := ParseToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }
        
        var user User
        if err := db.First(&user, userID).Error; err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
            c.Abort()
            return
        }
        
        if !user.IsActive {
            c.JSON(http.StatusForbidden, gin.H{"error": "account is disabled"})
            c.Abort()
            return
        }
        
        c.Set("user_id", userID)
        c.Set("user", user)
        c.Next()
    }
}