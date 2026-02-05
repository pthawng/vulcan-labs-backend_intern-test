package promotion

import "promotion-validator/internal/repository"

// PromotionService handles promotion eligibility use-cases.
type PromotionService struct {
	campaignRepo   repository.CodeRepository
	membershipRepo repository.CodeRepository
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
