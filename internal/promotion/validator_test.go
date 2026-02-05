package promotion

import (
	"testing"

	"promotion-validator/internal/repository"
)

func TestIsEligible_ExistsInBoth(t *testing.T) {
	campaignCodes := []string{"promo", "sale", "xyz"}
	membershipCodes := []string{"promo", "gold"}

	campaignRepo := repository.NewMockCodeRepository(campaignCodes, nil)
	membershipRepo := repository.NewMockCodeRepository(membershipCodes, nil)
	service := NewPromotionService(campaignRepo, membershipRepo)

	result, err := service.IsEligible("promo")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != true {
		t.Errorf("expected true, got %v", result)
	}
}

func TestIsEligible_ExistsInOnlyOne(t *testing.T) {
	campaignCodes := []string{"promo", "sale", "xyz"}
	membershipCodes := []string{"promo", "gold"}

	campaignRepo := repository.NewMockCodeRepository(campaignCodes, nil)
	membershipRepo := repository.NewMockCodeRepository(membershipCodes, nil)
	service := NewPromotionService(campaignRepo, membershipRepo)

	// Test code only in campaign
	result, err := service.IsEligible("sale")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != false {
		t.Errorf("expected false for code only in campaign, got %v", result)
	}

	// Test code only in membership
	result, err = service.IsEligible("gold")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != false {
		t.Errorf("expected false for code only in membership, got %v", result)
	}
}

func TestIsEligible_ExistsInNeither(t *testing.T) {
	campaignCodes := []string{"promo", "sale"}
	membershipCodes := []string{"gold"}

	campaignRepo := repository.NewMockCodeRepository(campaignCodes, nil)
	membershipRepo := repository.NewMockCodeRepository(membershipCodes, nil)
	service := NewPromotionService(campaignRepo, membershipRepo)

	result, err := service.IsEligible("notfd")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != false {
		t.Errorf("expected false, got %v", result)
	}
}
