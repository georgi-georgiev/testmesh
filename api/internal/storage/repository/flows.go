package repository

import (
	"github.com/georgi-georgiev/testmesh/internal/storage/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FlowRepository handles flow database operations
type FlowRepository struct {
	db *gorm.DB
}

// NewFlowRepository creates a new flow repository
func NewFlowRepository(db *gorm.DB) *FlowRepository {
	return &FlowRepository{db: db}
}

// Create creates a new flow
func (r *FlowRepository) Create(flow *models.Flow) error {
	return r.db.Create(flow).Error
}

// GetByID retrieves a flow by ID
func (r *FlowRepository) GetByID(id uuid.UUID) (*models.Flow, error) {
	var flow models.Flow
	if err := r.db.First(&flow, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &flow, nil
}

// List retrieves flows with optional filters
func (r *FlowRepository) List(suite string, tags []string, limit, offset int) ([]models.Flow, int64, error) {
	var flows []models.Flow
	var total int64

	query := r.db.Model(&models.Flow{})

	// Apply filters
	if suite != "" {
		query = query.Where("suite = ?", suite)
	}
	if len(tags) > 0 {
		query = query.Where("tags && ?", tags)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&flows).Error; err != nil {
		return nil, 0, err
	}

	return flows, total, nil
}

// Update updates a flow
func (r *FlowRepository) Update(flow *models.Flow) error {
	return r.db.Save(flow).Error
}

// Delete deletes a flow (soft delete)
func (r *FlowRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Flow{}, "id = ?", id).Error
}
