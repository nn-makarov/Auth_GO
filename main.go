package main

import (
    "github.com/gin-gonic/gin"
)

func main() {
    // Инициализация БД
    db := InitDB()
    
    // Создаем роутер
    r := gin.Default()
    
    // Публичные маршруты
    r.POST("/register", RegisterHandler(db))
    r.POST("/login", LoginHandler(db))
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Auth system on Go"})
    })
    
    // Защищенные маршруты (требуют токен)
    auth := r.Group("/")
    auth.Use(AuthMiddleware(db))
    {
        auth.GET("/profile", ProfileHandler(db))
        auth.GET("/api/products", ProductsHandler(db))
        auth.POST("/admin/create-permission", CreatePermissionHandler(db))
    }
    
    // Запуск сервера
    r.Run(":8080")
}