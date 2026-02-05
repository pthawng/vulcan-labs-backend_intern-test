package main

import (
	"fmt"
	"os"

	"promotion-validator/internal/promotion"
	"promotion-validator/internal/repository"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("usage: app <code> <campaign_file> <membership_file>")
		os.Exit(1)
	}

	code := os.Args[1]
	campaignFile := os.Args[2]
	membershipFile := os.Args[3]

	campaignRepo := repository.NewFileCodeRepository(campaignFile)
	membershipRepo := repository.NewFileCodeRepository(membershipFile)

	service := promotion.NewPromotionService(campaignRepo, membershipRepo)

	ok, err := service.IsEligible(code)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Use fmt.Fprintf to explicitly write to stdout
	// This ensures output is visible in PowerShell
	fmt.Fprintf(os.Stdout, "%v\n", ok)
	os.Stdout.Sync() // Flush stdout buffer
}
