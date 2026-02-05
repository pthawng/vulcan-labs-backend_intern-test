package promotion

// loadCampaignSet loads campaign codes into memory exactly once using sync.Once.
// This ensures thread-safe lazy initialization without locks on subsequent calls.
func (s *PromotionService) loadCampaignSet() error {
	s.loadOnce.Do(func() {
		var set map[string]struct{}
		set, s.loadErr = s.campaignRepo.LoadAll()
		if s.loadErr == nil {
			s.campaignSet = set
		}
	})
	return s.loadErr
}

// IsEligible checks whether a promotion code exists in both systems.
//
// Algorithm:
// 1. Validate input (length, characters)
// 2. Load campaign codes into HashSet (once, cached with sync.Once)
// 3. Check if code exists in campaign set - O(1) lookup
// 4. If yes, stream membership file to verify - O(m) with early exit
//
// Performance:
// - First call: O(n) to load campaign + O(1) lookup + O(m) membership check
// - Subsequent calls: O(1) lookup + O(m) membership check
//
// Memory:
// - Campaign HashSet: ~20-30MB for 5M codes (loaded once)
// - Membership streaming: O(1) buffer (~4KB)
//
// Thread-safety:
// - sync.Once ensures campaign set is loaded exactly once
// - Safe for concurrent calls to IsEligible
func (s *PromotionService) IsEligible(code string) (bool, error) {
	// Validate input
	if err := ValidateCode(code); err != nil {
		return false, err
	}

	// Load campaign codes into memory (once)
	if err := s.loadCampaignSet(); err != nil {
		return false, err
	}

	// O(1) lookup in campaign set
	if _, exists := s.campaignSet[code]; !exists {
		return false, nil // Early exit if not in campaign
	}

	// Stream membership file to verify
	// This is O(m) but with early exit when code is found
	return s.membershipRepo.Exists(code)
}
