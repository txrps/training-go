package repository

import (
	"context"
	"errors"
	"time"
	"training-go/internal/models"

	"gorm.io/gorm"
)

func CreateTodoGorm(db *gorm.DB, title string, completed bool) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	todo := &models.Todo{
		Title:     title,
		Completed: completed,
	}

	if err := db.WithContext(ctx).Create(todo).Error; err != nil {
		return nil, err
	}

	return todo, nil
}

func GetAllTodosGorm(db *gorm.DB) ([]models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var todos []models.Todo

	if err := db.WithContext(ctx).
		Order("created_at DESC").
		Find(&todos).Error; err != nil {
		return nil, err
	}

	return todos, nil
}

var ErrTodoNotFound = errors.New("todo not found")

func GetTodoByIDGorm(db *gorm.DB, id int64) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var todo models.Todo

	err := db.WithContext(ctx).
		First(&todo, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTodoNotFound
		}
		return nil, err
	}

	return &todo, nil
}

func UpdateTodoGorm(db *gorm.DB, id int64, title *string, completed *bool) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updates := map[string]interface{}{}

	if title != nil {
		updates["title"] = *title
	}

	if completed != nil {
		updates["completed"] = *completed
	}

	if len(updates) == 0 {
		return nil, errors.New("no fields to update")
	}

	var todo models.Todo

	result := db.WithContext(ctx).
		Model(&todo).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrTodoNotFound
	}

	if err := db.WithContext(ctx).First(&todo, id).Error; err != nil {
		return nil, err
	}

	return &todo, nil
}

func DeleteTodoGorm(db *gorm.DB, id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := db.WithContext(ctx).
		Delete(&models.Todo{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrTodoNotFound
	}

	return nil
}
