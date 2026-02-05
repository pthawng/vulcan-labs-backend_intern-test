package promotion

// IsEligible checks whether a promotion code exists in both systems.
func (s *PromotionService) IsEligible(code string) (bool, error) {
	// Validate input
	if err := ValidateCode(code); err != nil {
		return false, err
	}

	existsInCampaign, err := s.campaignRepo.Exists(code)
	if err != nil || !existsInCampaign {
		return false, err
	}

	existsInMembership, err := s.membershipRepo.Exists(code)
	if err != nil {
		return false, err
	}

	return existsInMembership, nil
}
