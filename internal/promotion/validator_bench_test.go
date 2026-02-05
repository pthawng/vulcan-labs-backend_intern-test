package promotion

import (
	"testing"

	"promotion-validator/internal/repository"
)

func BenchmarkIsEligible_HappyPath(b *testing.B) {
	campaignCodes := make([]string, 1000)
	membershipCodes := make([]string, 1000)

	for i := 0; i < 1000; i++ {
		code := string(rune('a'+i%26)) + string(rune('a'+(i/26)%26))
		campaignCodes[i] = code
		membershipCodes[i] = code
	}

	campaignRepo := repository.NewMockCodeRepository(campaignCodes, nil)
	membershipRepo := repository.NewMockCodeRepository(membershipCodes, nil)
	service := NewPromotionService(campaignRepo, membershipRepo)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.IsEligible("aa")
	}
}
