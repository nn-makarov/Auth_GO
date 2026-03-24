package main

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
	gorm.Model
    ID           uint      `gorm:"primaryKey"`
    Email        string    `gorm:"uniqueIndex;not null"`
    PasswordHash string    `gorm:"not null"`
    FirstName    string
    LastName     string
    MiddleName   string
    IsActive     bool      `gorm:"default:true"`
    Role         string    `gorm:"default:user"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type Resource struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string `gorm:"uniqueIndex;not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Permission struct {
    ID         uint `gorm:"primaryKey"`
    Role       string `gorm:"index:idx_role_resource,unique"`
    ResourceID uint   `gorm:"index:idx_role_resource,unique"`
    Resource   Resource
    CanView    bool `gorm:"default:false"`
    CanCreate  bool `gorm:"default:false"`
    CanEdit    bool `gorm:"default:false"`
    CanDelete  bool `gorm:"default:false"`
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

type RegisterRequest struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=8"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Token string `json:"token"`
    Email string `json:"email"`
    Role  string `json:"role"`
}

type CreatePermissionRequest struct {
    Role       string `json:"role" binding:"required"`
    Resource   string `json:"resource" binding:"required"`
    CanView    bool   `json:"can_view"`
    CanCreate  bool   `json:"can_create"`
    CanEdit    bool   `json:"can_edit"`
    CanDelete  bool   `json:"can_delete"`
}