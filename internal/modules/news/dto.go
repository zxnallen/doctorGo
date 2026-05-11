package news

import "time"

type ListRequest struct {
	Page       int    `form:"page"`
	Size       int    `form:"size"`
	CategoryID uint64 `form:"category_id"`
	Keyword    string `form:"keyword"`
}

type AdminListRequest struct {
	Page       int    `form:"page"`
	Size       int    `form:"size"`
	CategoryID uint64 `form:"category_id"`
	Status     *int   `form:"status"`
	Keyword    string `form:"keyword"`
}

type ListResponse struct {
	Items []Item `json:"items"`
	Total int64  `json:"total"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}

type Item struct {
	ID          uint64     `json:"id"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	CoverURL    string     `json:"cover_url"`
	Author      string     `json:"author"`
	CategoryID  uint64     `json:"category_id"`
	Status      int        `json:"status"`
	ViewCount   int64      `json:"view_count"`
	PublishedAt *time.Time `json:"published_at"`
}

type Detail struct {
	ID          uint64     `json:"id"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	Content     string     `json:"content"`
	CoverURL    string     `json:"cover_url"`
	Author      string     `json:"author"`
	CategoryID  uint64     `json:"category_id"`
	ViewCount   int64      `json:"view_count"`
	PublishedAt *time.Time `json:"published_at"`
}

type CreateRequest struct {
	Title       string     `json:"title" binding:"required,max=200"`
	Summary     string     `json:"summary" binding:"omitempty,max=500"`
	Content     string     `json:"content" binding:"omitempty"`
	CoverURL    string     `json:"cover_url" binding:"omitempty,max=512"`
	Author      string     `json:"author" binding:"omitempty,max=64"`
	CategoryID  uint64     `json:"category_id"`
	Status      int        `json:"status" binding:"oneof=0 1 2"`
	PublishedAt *time.Time `json:"published_at"`
}

type UpdateRequest struct {
	Title       string     `json:"title" binding:"required,max=200"`
	Summary     string     `json:"summary" binding:"omitempty,max=500"`
	Content     string     `json:"content" binding:"omitempty"`
	CoverURL    string     `json:"cover_url" binding:"omitempty,max=512"`
	Author      string     `json:"author" binding:"omitempty,max=64"`
	CategoryID  uint64     `json:"category_id"`
	Status      int        `json:"status" binding:"oneof=0 1 2"`
	PublishedAt *time.Time `json:"published_at"`
}

type UpdateStatusRequest struct {
	Status int `json:"status" binding:"oneof=0 1 2"`
}

type CategoryListRequest struct {
	OnlyEnabled bool `form:"only_enabled"`
}

type CategoryItem struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
}

type CreateCategoryRequest struct {
	Name   string `json:"name" binding:"required,max=64"`
	Sort   int    `json:"sort"`
	Status int    `json:"status" binding:"oneof=0 1"`
}

type UpdateCategoryRequest struct {
	Name   string `json:"name" binding:"required,max=64"`
	Sort   int    `json:"sort"`
	Status int    `json:"status" binding:"oneof=0 1"`
}

type UpdateCategoryStatusRequest struct {
	Status int `json:"status" binding:"oneof=0 1"`
}
