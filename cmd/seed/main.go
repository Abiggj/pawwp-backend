package main

import (
	"log"

	"github.com/Abiggj/pawwp/internal/account"
	"github.com/Abiggj/pawwp/internal/database"
	"github.com/Abiggj/pawwp/internal/pet"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()
	database.Connect()

	// 1. Clear existing data (Optional, for fresh seed)
	log.Println("Clearing existing data...")
	database.DB.Exec("TRUNCATE accounts, pets CASCADE")

	// 2. Create Accounts
	pass, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashedPass := string(pass)

	accounts := []account.Account{
		{
			Type:         "individual",
			Email:        "john@example.com",
			PasswordHash: hashedPass,
			Bio:          "Animal lover from NYC",
			IsVerified:   true,
		},
		{
			Type:         "shelter",
			Email:        "happy_paws@shelter.org",
			PasswordHash: hashedPass,
			Bio:          "Local rescue shelter",
			IsVerified:   true,
		},
	}

	for i := range accounts {
		if err := database.DB.Create(&accounts[i]).Error; err != nil {
			log.Fatalf("Failed to seed account: %v", err)
		}
	}

	// 3. Create Pets
	pets := []pet.Pet{
		{
			AccountID:      accounts[0].ID,
			Username:       "buddy_the_golden",
			Name:           "Buddy",
			Species:        "Dog",
			Breed:          "Golden Retriever",
			Bio:            "I love tennis balls!",
			AdoptionStatus: "not_applicable",
		},
		{
			AccountID:      accounts[1].ID,
			Username:       "mittens_rescue",
			Name:           "Mittens",
			Species:        "Cat",
			Breed:          "Tabby",
			Bio:            "Looking for a forever home.",
			AdoptionStatus: "adoptable",
		},
	}

	for i := range pets {
		if err := database.DB.Create(&pets[i]).Error; err != nil {
			log.Fatalf("Failed to seed pet: %v", err)
		}
	}

	log.Println("✅ Database successfully seeded!")
}
