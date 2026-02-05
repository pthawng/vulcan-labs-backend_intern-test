package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	totalCampaign = 500_000
	overlapRatio  = 0.4
)

func randomCode() string {
	length := rand.Intn(5) + 1
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = byte('a' + rand.Intn(26))
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	_ = os.MkdirAll("data", 0755)

	campaignFile, err := os.Create("data/campaign_codes.txt")
	if err != nil {
		panic(err)
	}
	defer campaignFile.Close()

	membershipFile, err := os.Create("data/membership_codes.txt")
	if err != nil {
		panic(err)
	}
	defer membershipFile.Close()

	campaignSet := make(map[string]struct{})
	membershipSet := make(map[string]struct{})

	// 1️⃣ Generate campaign codes
	for len(campaignSet) < totalCampaign {
		code := randomCode()
		campaignSet[code] = struct{}{}
	}

	// 2️⃣ Write campaign & probabilistic overlap
	for code := range campaignSet {
		fmt.Fprintln(campaignFile, code)
		if rand.Float64() < overlapRatio {
			membershipSet[code] = struct{}{}
		}
	}

	// 3️⃣ Add membership-only codes
	targetMembership := int(float64(totalCampaign) * 0.6)
	for len(membershipSet) < targetMembership {
		code := randomCode()
		if _, exists := campaignSet[code]; !exists {
			membershipSet[code] = struct{}{}
		}
	}

	// 4️⃣ Write membership codes
	for code := range membershipSet {
		fmt.Fprintln(membershipFile, code)
	}

	fmt.Println("✅ Dataset generated (pure random, matches problem constraints)")
	fmt.Printf("Campaign codes: %d\n", len(campaignSet))
	fmt.Printf("Membership codes: %d\n", len(membershipSet))
}
