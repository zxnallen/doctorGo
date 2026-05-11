package news

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"doctor-go/internal/infrastructure/redis"
	"doctor-go/internal/pkg/pagination"
)

const (
	cacheKeyNewsListPattern      = "cache:news:list:*"
	cacheKeyNewsDetailPrefix     = "cache:news:detail:"
	cacheKeyNewsCategoriesPrefix = "cache:news:categories:"
)

var ErrNewsNotFound = errors.New("news not found")
var ErrCategoryNotFound = errors.New("news category not found")

type Service struct {
	repo  *Repository
	cache *redis.Client
}

func NewService(repo *Repository, cache *redis.Client) *Service {
	return &Service{repo: repo, cache: cache}
}

func (s *Service) List(ctx context.Context, req ListRequest) (*ListResponse, error) {
	page := pagination.Page{Page: req.Page, Size: req.Size}.Normalize()
	cacheKey := fmt.Sprintf("cache:news:list:page=%d:size=%d:category=%d:keyword=%s", page.Page, page.Size, req.CategoryID, req.Keyword)
	var cached ListResponse
	if s.getCache(ctx, cacheKey, &cached) {
		return &cached, nil
	}

	items, total, err := s.repo.List(req, page.Offset(), page.Size)
	if err != nil {
		return nil, err
	}

	result := &ListResponse{
		Items: toItems(items),
		Total: total,
		Page:  page.Page,
		Size:  page.Size,
	}
	s.setCache(ctx, cacheKey, result, 5*time.Minute)
	return result, nil
}

func (s *Service) AdminList(req AdminListRequest) (*ListResponse, error) {
	page := pagination.Page{Page: req.Page, Size: req.Size}.Normalize()
	items, total, err := s.repo.AdminList(req, page.Offset(), page.Size)
	if err != nil {
		return nil, err
	}

	return &ListResponse{
		Items: toItems(items),
		Total: total,
		Page:  page.Page,
		Size:  page.Size,
	}, nil
}

func (s *Service) Detail(ctx context.Context, id uint64) (*Detail, error) {
	cacheKey := cacheKeyNewsDetailPrefix + strconv.FormatUint(id, 10)
	var cached Detail
	if s.getCache(ctx, cacheKey, &cached) {
		s.incrementViewCount(ctx, id)
		cached.ViewCount++
		return &cached, nil
	}

	item, err := s.repo.FindPublishedByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrNewsNotFound
	}

	s.incrementViewCount(ctx, id)

	detail := toDetail(*item)
	detail.ViewCount++
	s.setCache(ctx, cacheKey, &detail, 10*time.Minute)
	return &detail, nil
}

func (s *Service) AdminDetail(id uint64) (*Detail, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrNewsNotFound
	}
	detail := toDetail(*item)
	return &detail, nil
}

func (s *Service) Create(req CreateRequest) (*Detail, error) {
	item := &News{
		Title:       req.Title,
		Summary:     req.Summary,
		Content:     req.Content,
		CoverURL:    req.CoverURL,
		Author:      req.Author,
		CategoryID:  req.CategoryID,
		Status:      req.Status,
		PublishedAt: publishedAt(req.Status, req.PublishedAt),
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	s.invalidateNewsCache(context.Background(), item.ID)
	detail := toDetail(*item)
	return &detail, nil
}

func (s *Service) Update(id uint64, req UpdateRequest) (*Detail, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrNewsNotFound
	}

	item.Title = req.Title
	item.Summary = req.Summary
	item.Content = req.Content
	item.CoverURL = req.CoverURL
	item.Author = req.Author
	item.CategoryID = req.CategoryID
	item.Status = req.Status
	item.PublishedAt = publishedAt(req.Status, req.PublishedAt)

	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	s.invalidateNewsCache(context.Background(), item.ID)
	detail := toDetail(*item)
	return &detail, nil
}

func (s *Service) UpdateStatus(id uint64, status int) (*Detail, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrNewsNotFound
	}

	item.Status = status
	item.PublishedAt = publishedAt(status, item.PublishedAt)
	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	s.invalidateNewsCache(context.Background(), item.ID)
	detail := toDetail(*item)
	return &detail, nil
}

func (s *Service) Delete(id uint64) error {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if item == nil {
		return ErrNewsNotFound
	}
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	s.invalidateNewsCache(context.Background(), id)
	return nil
}

func (s *Service) ListCategories(ctx context.Context, req CategoryListRequest) ([]CategoryItem, error) {
	cacheKey := cacheKeyNewsCategoriesPrefix + strconv.FormatBool(req.OnlyEnabled)
	var cached []CategoryItem
	if s.getCache(ctx, cacheKey, &cached) {
		return cached, nil
	}

	items, err := s.repo.ListCategories(req.OnlyEnabled)
	if err != nil {
		return nil, err
	}
	result := toCategoryItems(items)
	s.setCache(ctx, cacheKey, result, 30*time.Minute)
	return result, nil
}

func (s *Service) CreateCategory(req CreateCategoryRequest) (*CategoryItem, error) {
	item := &Category{
		Name:   req.Name,
		Sort:   req.Sort,
		Status: req.Status,
	}
	if err := s.repo.CreateCategory(item); err != nil {
		return nil, err
	}
	s.invalidateCategoryCache(context.Background())
	result := toCategoryItem(*item)
	return &result, nil
}

func (s *Service) UpdateCategory(id uint64, req UpdateCategoryRequest) (*CategoryItem, error) {
	item, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrCategoryNotFound
	}

	item.Name = req.Name
	item.Sort = req.Sort
	item.Status = req.Status
	if err := s.repo.UpdateCategory(item); err != nil {
		return nil, err
	}
	s.invalidateNewsCache(context.Background(), 0)
	s.invalidateCategoryCache(context.Background())
	result := toCategoryItem(*item)
	return &result, nil
}

func (s *Service) UpdateCategoryStatus(id uint64, status int) (*CategoryItem, error) {
	item, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrCategoryNotFound
	}

	item.Status = status
	if err := s.repo.UpdateCategory(item); err != nil {
		return nil, err
	}
	s.invalidateNewsCache(context.Background(), 0)
	s.invalidateCategoryCache(context.Background())
	result := toCategoryItem(*item)
	return &result, nil
}

func (s *Service) DeleteCategory(id uint64) error {
	item, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return err
	}
	if item == nil {
		return ErrCategoryNotFound
	}
	if err := s.repo.DeleteCategory(id); err != nil {
		return err
	}
	s.invalidateNewsCache(context.Background(), 0)
	s.invalidateCategoryCache(context.Background())
	return nil
}

func (s *Service) incrementViewCount(ctx context.Context, id uint64) {
	if s.cache != nil {
		_ = s.cache.Incr(ctx, "news:view_count:"+strconv.FormatUint(id, 10)).Err()
	}
	_ = s.repo.IncrementViewCount(id)
}

func (s *Service) getCache(ctx context.Context, key string, target interface{}) bool {
	if s.cache == nil {
		return false
	}
	value, err := s.cache.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	return json.Unmarshal([]byte(value), target) == nil
}

func (s *Service) setCache(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	if s.cache == nil {
		return
	}
	bytes, err := json.Marshal(value)
	if err != nil {
		return
	}
	_ = s.cache.Set(ctx, key, string(bytes), ttl).Err()
}

func (s *Service) invalidateNewsCache(ctx context.Context, detailID uint64) {
	if s.cache == nil {
		return
	}
	s.deleteByPattern(ctx, cacheKeyNewsListPattern)
	if detailID > 0 {
		_ = s.cache.Del(ctx, cacheKeyNewsDetailPrefix+strconv.FormatUint(detailID, 10)).Err()
	}
}

func (s *Service) invalidateCategoryCache(ctx context.Context) {
	if s.cache == nil {
		return
	}
	s.deleteByPattern(ctx, cacheKeyNewsCategoriesPrefix+"*")
}

func (s *Service) deleteByPattern(ctx context.Context, pattern string) {
	iter := s.cache.Scan(ctx, 0, pattern, 100).Iterator()
	for iter.Next(ctx) {
		_ = s.cache.Del(ctx, iter.Val()).Err()
	}
}

func toItems(items []News) []Item {
	result := make([]Item, 0, len(items))
	for _, item := range items {
		result = append(result, Item{
			ID:          item.ID,
			Title:       item.Title,
			Summary:     item.Summary,
			CoverURL:    item.CoverURL,
			Author:      item.Author,
			CategoryID:  item.CategoryID,
			Status:      item.Status,
			ViewCount:   item.ViewCount,
			PublishedAt: item.PublishedAt,
		})
	}
	return result
}

func toDetail(item News) Detail {
	return Detail{
		ID:          item.ID,
		Title:       item.Title,
		Summary:     item.Summary,
		Content:     item.Content,
		CoverURL:    item.CoverURL,
		Author:      item.Author,
		CategoryID:  item.CategoryID,
		ViewCount:   item.ViewCount,
		PublishedAt: item.PublishedAt,
	}
}

func toCategoryItems(items []Category) []CategoryItem {
	result := make([]CategoryItem, 0, len(items))
	for _, item := range items {
		result = append(result, toCategoryItem(item))
	}
	return result
}

func toCategoryItem(item Category) CategoryItem {
	return CategoryItem{
		ID:     item.ID,
		Name:   item.Name,
		Sort:   item.Sort,
		Status: item.Status,
	}
}

func publishedAt(status int, current *time.Time) *time.Time {
	if status != 1 {
		return current
	}
	if current != nil {
		return current
	}
	now := time.Now()
	return &now
}
