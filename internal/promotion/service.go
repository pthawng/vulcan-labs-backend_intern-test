package promotion

import (
	"sync"

	"promotion-validator/internal/repository"
)

// PromotionService handles promotion eligibility use-cases.
type PromotionService struct {
	campaignRepo   repository.CodeRepository
	membershipRepo repository.CodeRepository

	// Campaign codes are loaded once and cached for O(1) lookup
	campaignSet map[string]struct{}
	loadOnce    sync.Once
	loadErr     error
}

func NewPromotionService(
	campaignRepo repository.CodeRepository,
	membershipRepo repository.CodeRepository,
) *PromotionService {
	return &PromotionService{
		campaignRepo:   campaignRepo,
		membershipRepo: membershipRepo,
	}
}
