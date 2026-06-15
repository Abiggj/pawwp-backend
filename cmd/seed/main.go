package main

import (
	"log"
	"time"

	"github.com/Abiggj/pawwp/internal/account"
	"github.com/Abiggj/pawwp/internal/community"
	"github.com/Abiggj/pawwp/internal/database"
	"github.com/Abiggj/pawwp/internal/pet"
	"github.com/Abiggj/pawwp/internal/post"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()
	database.Connect()

	// 1. Clear existing data
	log.Println("🧹 Clearing existing data...")
	database.DB.Exec("TRUNCATE accounts, pets, posts, boops, woofs, channels, channel_posts, subscriptions CASCADE")

	// 2. Create Accounts
	log.Println("🌱 Seeding Accounts...")
	pass, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashedPass := string(pass)

	parentAcc := account.Account{
		ID:           uuid.New(),
		Type:         "parent",
		Email:        "parent@example.com",
		PasswordHash: hashedPass,
		Bio:          "Proud dog parent and animal advocate.",
		IsVerified:   true,
	}

	shelterAcc := account.Account{
		ID:           uuid.New(),
		Type:         "shelter",
		Email:        "contact@happypaws.org",
		PasswordHash: hashedPass,
		Bio:          "Happy Paws Rescue - Dedicated to finding homes for every pet.",
		IsVerified:   true,
	}

	database.DB.Create(&parentAcc)
	database.DB.Create(&shelterAcc)

	// 3. Create Pets
	log.Println("🐕 Seeding Pets...")
	buddy := pet.Pet{
		ID:             uuid.New(),
		AccountID:      parentAcc.ID,
		Username:       "buddy_goldie",
		Name:           "Buddy",
		Species:        "Dog",
		Breed:          "Golden Retriever",
		Bio:            "Professional fetcher and treat connoisseur.",
		AdoptionStatus: "not_applicable",
	}

	luna := pet.Pet{
		ID:             uuid.New(),
		AccountID:      shelterAcc.ID,
		Username:       "luna_rescue",
		Name:           "Luna",
		Species:        "Cat",
		Breed:          "Siamese Mix",
		Bio:            "Quiet and loving, needs a calm home.",
		AdoptionStatus: "adoptable",
	}

	database.DB.Create(&buddy)
	database.DB.Create(&luna)

	// 4. Create Personal Media Posts
	log.Println("📸 Seeding Media Posts...")
	posts := []post.Post{
		{
			ID:          uuid.New(),
			PetID:       buddy.ID,
			MediaURL:    "https://images.unsplash.com/photo-1552053831-71594a27632d",
			Caption:     "Just got a new ball! #doglife",
			IsShowcased: true,
			CreatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			PetID:       buddy.ID,
			MediaURL:    "https://images.unsplash.com/photo-1583511655857-d19b40a7a54e",
			Caption:     "Sunny day at the park.",
			IsShowcased: false,
			CreatedAt:   time.Now().Add(-1 * time.Hour),
		},
		{
			ID:          uuid.New(),
			PetID:       luna.ID,
			MediaURL:    "https://images.unsplash.com/photo-1514888286974-6c03e2ca1dba",
			Caption:     "Waiting for my forever family.",
			IsShowcased: true,
			CreatedAt:   time.Now().Add(-5 * time.Hour),
		},
	}

	for _, p := range posts {
		database.DB.Create(&p)
	}

	// 5. Create Community Channels
	log.Println("🏘 Seeding Channels...")
	rescueChannel := community.Channel{
		ID:           uuid.New(),
		AccountID:    shelterAcc.ID,
		Name:         "Happy Paws Rescue Hub",
		Description:  "Emergency rescues and adoption calls for the tri-state area.",
		CoverageArea: "New York, New Jersey, Connecticut",
		IsActive:     true,
	}
	database.DB.Create(&rescueChannel)

	// 6. Subscriptions
	log.Println("🔔 Seeding Subscriptions...")
	sub := community.Subscription{
		AccountID: parentAcc.ID,
		ChannelID: rescueChannel.ID,
		CreatedAt: time.Now(),
	}
	database.DB.Create(&sub)

	// 7. Channel Posts (SOS, Adoption, Donation)
	log.Println("📢 Seeding Channel Posts...")
	cPosts := []community.ChannelPost{
		{
			ID:        uuid.New(),
			ChannelID: rescueChannel.ID,
			AuthorID:  shelterAcc.ID,
			Type:      "sos",
			Status:    "urgent",
			Title:     "Injured stray near 5th Ave",
			Content:   "Injured kitten found near 5th Ave. We need a volunteer with a carrier immediately!",
			MediaURL:  "https://images.unsplash.com/photo-1548546738-8502ce9a3e0c",
		},
		{
			ID:        uuid.New(),
			ChannelID: rescueChannel.ID,
			AuthorID:  shelterAcc.ID,
			Type:      "adoption",
			Title:     "Meet Luna - Ready for Adoption",
			Content:   "Luna is a sweet 2-year-old Siamese mix. She is vaccinated and spayed.",
			MediaURL:  "https://images.unsplash.com/photo-1514888286974-6c03e2ca1dba",
		},
		{
			ID:        uuid.New(),
			ChannelID: rescueChannel.ID,
			AuthorID:  shelterAcc.ID,
			Type:      "donation",
			Title:     "Winter Blanket Drive",
			Content:   "Winter is coming! We are raising $500 for heavy-duty blankets for our outdoor kennels.",
		},
	}

	for _, cp := range cPosts {
		database.DB.Create(&cp)
	}

	log.Println("✨ Database successfully seeded with Pawwp data!")
}
