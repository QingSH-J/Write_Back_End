package storage

//imports
import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinxinyu/go_backend/internal/models"
	"gorm.io/gorm"
)

// WriteLogRepository defines the interface for write log operations
type WriteLogRepository interface {
	CreateLog(ctx context.Context, log *models.WriteLog) error
	UpdateLog(ctx context.Context, log *models.WriteLog) error
	GetLogByID(ctx context.Context, id uuid.UUID) (*models.WriteLog, error)
	GetLogByUserID(ctx context.Context, userID uuid.UUID) ([]*models.WriteLog, error)
	GetLogByDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]*models.WriteLog, error)
	GetLogByDateRange(ctx context.Context, userID uuid.UUID, startDate time.Time, endDate time.Time) ([]*models.WriteLog, error)
	DeleteLog(ctx context.Context, id uuid.UUID) error
}

type writeLogRepository struct {
	db *gorm.DB
}

func NewWriteLogRepository(db *gorm.DB) WriteLogRepository {
	return &writeLogRepository{db: db}
}

func (r *writeLogRepository) CreateLog(ctx context.Context, log *models.WriteLog) error {
	result := r.db.WithContext(ctx).Create(log)
	if result.Error != nil {
		return fmt.Errorf("failed to create log: %v", result.Error)
	}
	return nil
}

func (r *writeLogRepository) UpdateLog(ctx context.Context, log *models.WriteLog) error {
	result := r.db.WithContext(ctx).Model(&models.WriteLog{}).Where("id = ? AND user_id = ?", log.ID, log.UserID).Updates(map[string]interface{}{
		"word_count": log.WordsCount,
		"content":    log.Content,
		//gorm will automatically update the updated_at field because of the autoUpdateTime
	})
	//check if the log was updated
	if result.RowsAffected == 0 {
		return fmt.Errorf("log not found(ID: %s, UserID: %s)", log.ID, log.UserID)
	}
	return nil
}

func (r *writeLogRepository) GetLogByID(ctx context.Context, id uuid.UUID) (*models.WriteLog, error) {
	var log models.WriteLog
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&log)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("log not found with id %s", id)
		}
		return nil, fmt.Errorf("failed to get log: %w", result.Error)
	}
	return &log, nil
}

func (r *writeLogRepository) GetLogByUserID(ctx context.Context, userID uuid.UUID) ([]*models.WriteLog, error) {
	var logs []*models.WriteLog
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&logs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get logs: %v", result.Error)
	}
	return logs, nil
}

func (r *writeLogRepository) GetLogByDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]*models.WriteLog, error) {
	var logs []*models.WriteLog
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	result := r.db.WithContext(ctx).Where("user_id = ? AND date BETWEEN ? AND ?", userID, startOfDay, endOfDay).Order("created_at asc").Find(&logs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get logs: %v", result.Error)
	}
	return logs, nil
}

func (r *writeLogRepository) GetLogByDateRange(ctx context.Context, userID uuid.UUID, startDate time.Time, endDate time.Time) ([]*models.WriteLog, error) {
	var logs []*models.WriteLog
	result := r.db.WithContext(ctx).Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).Order("created_at asc").Find(&logs)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get logs: %v", result.Error)
	}
	return logs, nil
}

func (r *writeLogRepository) DeleteLog(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.WriteLog{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete log: %v", result.Error)
	}
	return nil
}
