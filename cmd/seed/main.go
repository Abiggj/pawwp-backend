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
	log.Println("Clearing existing data...")
	database.DB.Exec("TRUNCATE accounts, pets, posts, boops, woofs, channels, channel_posts, subscriptions, channel_admins, townhall_questions, townhall_answers CASCADE")

	// 2. Create Accounts
	log.Println("Seeding Diverse Accounts...")
	pass, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	hashedPass := string(pass)

	// Pet Parents
	parentJohn := account.Account{
		ID:           uuid.New(),
		Type:         "parent",
		Email:        "john@example.com",
		PasswordHash: hashedPass,
		Bio:          "Dog dad to Buddy and Milo. Passionate about weekend hiking and agility trails in Central Park.",
		IsVerified:   true,
		CreatedAt:    time.Now().Add(-60 * 24 * time.Hour),
	}

	parentElena := account.Account{
		ID:           uuid.New(),
		Type:         "parent",
		Email:        "elena@example.com",
		PasswordHash: hashedPass,
		Bio:          "Cat mom to Cleo. Fostering rescue kittens and advocating for animal welfare in Brooklyn.",
		IsVerified:   true,
		CreatedAt:    time.Now().Add(-45 * 24 * time.Hour),
	}

	parentMarcus := account.Account{
		ID:           uuid.New(),
		Type:         "parent",
		Email:        "marcus@example.com",
		PasswordHash: hashedPass,
		Bio:          "Rocky the Frenchie's companion. Coffee and canine enthusiast in Queens.",
		IsVerified:   false,
		CreatedAt:    time.Now().Add(-20 * 24 * time.Hour),
	}

	// Shelters & NGOs
	shelterHappyPaws := account.Account{
		ID:           uuid.New(),
		Type:         "shelter",
		Email:        "contact@happypaws.org",
		PasswordHash: hashedPass,
		Bio:          "Happy Paws Rescue Hub (501c3). Dedicated to rescue, rehabilitation, and adoption across the Tri-State area.",
		IsVerified:   true,
		CreatedAt:    time.Now().Add(-180 * 24 * time.Hour),
	}

	shelterCityStrays := account.Account{
		ID:           uuid.New(),
		Type:         "shelter",
		Email:        "rescue@citystrays.org",
		PasswordHash: hashedPass,
		Bio:          "City Stray Animal Alliance. Rapid emergency dispatch for injured and abandoned urban animals.",
		IsVerified:   true,
		CreatedAt:    time.Now().Add(-120 * 24 * time.Hour),
	}

	shelterSecondChance := account.Account{
		ID:           uuid.New(),
		Type:         "shelter",
		Email:        "haven@secondchance.org",
		PasswordHash: hashedPass,
		Bio:          "Second Chance Animal Sanctuary. Providing lifelong care for special needs and senior rescues.",
		IsVerified:   true,
		CreatedAt:    time.Now().Add(-90 * 24 * time.Hour),
	}

	accounts := []account.Account{
		parentJohn,
		parentElena,
		parentMarcus,
		shelterHappyPaws,
		shelterCityStrays,
		shelterSecondChance,
	}

	for i := range accounts {
		database.DB.Create(&accounts[i])
	}

	// 3. Create Pets
	log.Println("Seeding Diverse Pet Profiles...")

	buddy := pet.Pet{
		ID:             uuid.New(),
		AccountID:      parentJohn.ID,
		Username:       "buddy_goldie",
		Name:           "Buddy",
		Species:        "Dog",
		Breed:          "Golden Retriever",
		Bio:            "Professional fetcher, water lover, and treat connoisseur.",
		AdoptionStatus: "not_applicable",
		IsVerified:     true,
	}

	milo := pet.Pet{
		ID:             uuid.New(),
		AccountID:      parentJohn.ID,
		Username:       "milo_corgi",
		Name:           "Milo",
		Species:        "Dog",
		Breed:          "Pembroke Welsh Corgi",
		Bio:            "Short legs, massive personality. Squirrel patrol specialist.",
		AdoptionStatus: "not_applicable",
		IsVerified:     false,
	}

	cleo := pet.Pet{
		ID:             uuid.New(),
		AccountID:      parentElena.ID,
		Username:       "cleo_queen",
		Name:           "Cleo",
		Species:        "Cat",
		Breed:          "Persian",
		Bio:            "Queen of sunbeams and high shelves. Tolerates petting on her terms.",
		AdoptionStatus: "not_applicable",
		IsVerified:     true,
	}

	rocky := pet.Pet{
		ID:             uuid.New(),
		AccountID:      parentMarcus.ID,
		Username:       "rocky_frenchie",
		Name:           "Rocky",
		Species:        "Dog",
		Breed:          "French Bulldog",
		Bio:            "Expert snorer, hoodie model, and neighborhood celebrity.",
		AdoptionStatus: "not_applicable",
		IsVerified:     false,
	}

	luna := pet.Pet{
		ID:             uuid.New(),
		AccountID:      shelterHappyPaws.ID,
		Username:       "luna_rescue",
		Name:           "Luna",
		Species:        "Cat",
		Breed:          "Siamese Mix",
		Bio:            "Gentle 2yo Siamese mix. Fully vaccinated, spayed, and looking for a quiet home.",
		AdoptionStatus: "adoptable",
		IsVerified:     true,
	}

	barnaby := pet.Pet{
		ID:             uuid.New(),
		AccountID:      shelterHappyPaws.ID,
		Username:       "barnaby_beagle",
		Name:           "Barnaby",
		Species:        "Dog",
		Breed:          "Beagle Mix",
		Bio:            "Sweet 1yo hound dog with floppy ears. Great with kids and other companions.",
		AdoptionStatus: "adoptable",
		IsVerified:     true,
	}

	oliver := pet.Pet{
		ID:             uuid.New(),
		AccountID:      shelterCityStrays.ID,
		Username:       "oliver_tabby",
		Name:           "Oliver",
		Species:        "Cat",
		Breed:          "Orange Tabby",
		Bio:            "Rescued street kitten! Very playful purr-machine seeking forever family.",
		AdoptionStatus: "adoptable",
		IsVerified:     true,
	}

	bruno := pet.Pet{
		ID:             uuid.New(),
		AccountID:      shelterSecondChance.ID,
		Username:       "bruno_shepherd",
		Name:           "Bruno",
		Species:        "Dog",
		Breed:          "German Shepherd",
		Bio:            "Loyal senior shepherd rescued from neglect. Thriving in sanctuary care.",
		AdoptionStatus: "rescued",
		IsVerified:     true,
	}

	pets := []pet.Pet{buddy, milo, cleo, rocky, luna, barnaby, oliver, bruno}
	for i := range pets {
		database.DB.Create(&pets[i])
	}

	// 4. Create Posts (Top Catches & Showcase Items)
	log.Println("Seeding Media Posts and 9-Post Showcases...")

	posts := []post.Post{
		// Buddy's Posts (Complete 9-Post Showcase Reel)
		{
			ID:          uuid.New(),
			PetID:       buddy.ID,
			MediaURL:    "https://images.unsplash.com/photo-1552053831-71594a27632d?auto=format&fit=crop&w=800&q=80",
			Caption:     "Found the biggest stick in Central Park today. Absolutely refusing to drop it.",
			IsShowcased: true,
			CreatedAt:   time.Now().Add(-2 * time.Hour),
		},
		{
			ID:          uuid.New(),
			PetID:       buddy.ID,
			MediaURL:    "https://images.unsplash.com/photo-1583511655857-d19b40a7a54e?auto=format&fit=crop&w=800&q=80",
			Caption:     "Sunny morning park adventures with the squad.",
			IsShowcased: true,
			CreatedAt:   time.Now().Add(-6 * time.Hour),
		},
		{
			ID:          uuid.New(),
			PetID:       buddy.ID,
			MediaURL:    "https://images.unsplash.com/photo-1517849845537-4d257902454a?auto=format&fit=crop&w=800&q=80",
			Caption:     "Post-bath zoomies incoming!",
			IsShowcased: true,
			CreatedAt:   time.Now().Add(-18 * time.Hour),
		},
		{
			ID:          uuid.New(),
			PetID:       buddy.ID,
			MediaURL:    "https://images.unsplash.com/photo-1537151625747-768eb6cf92b2?auto=format&fit=crop&w=800&q=80",
			Caption:     "Waiting patiently by the treat cabinet. It worked!",
			IsShowcased: true,
			CreatedAt:   time.Now().Add(-28 * time.Hour),
		},
		{
			ID:          uuid.New(),
			PetID:       buddy.ID,
			MediaURL:    "https://images.unsplash.com/photo-1561037404-61cd46aa615b?auto=format&fit=crop&w=800&q=80",
			Caption:     "First time seeing snow this year! Thoroughly fascinated.",
			IsShowcased: true,
			CreatedAt:   time.Now().Add(-35 * time.Hour),
		},
		{
			ID:          uuid.New(),
			PetID:       buddy.ID,
			MediaURL:    "https://images.unsplash.com/photo-1543466835-00a7907e9de1?auto=format&fit=crop&w=800&q=80",
			Caption:     "Head tilts whenever peanut butter is mentioned.",
			IsShowcased: true,
			CreatedAt:   time.Now().Add(-40 * time.Hour),
		},
		{
			ID:          uuid.New(),
			PetID:       buddy.ID,
			MediaURL:    "https://images.unsplash.com/photo-1587300003388-59208cc962cb?auto=format&fit=crop&w=800&q=80",
			Caption:     "Beach day retriever mode activated.",
			IsShowcased: true,
			CreatedAt:   time.Now().Add(-44 * time.Hour),
		},
		{
			ID:          uuid.New(),
			PetID:       buddy.ID,
			MediaURL:    "https://images.unsplash.com/photo-1598133894008-61f7fdb8cc3a?auto=format&fit=crop&w=800&q=80",
			Caption:     "Afternoon snooze in the sunbeam. Do not disturb.",
			IsShowcased: true,
			CreatedAt:   time.Now().Add(-46 * time.Hour),
		},
		{
			ID:          uuid.New(),
			PetID:       buddy.ID,
			MediaURL:    "https://images.unsplash.com/photo-1548767797-d8c844163c4c?auto=format&fit=crop&w=800&q=80",
			Caption:     "The legendary puppy throwback photo! Where did the time go?",
			IsShowcased: true,
			CreatedAt:   time.Now().Add(-47 * time.Hour),
		},

		// Luna's Posts (Adoptable Siamese Mix)
		{
			ID:          uuid.New(),
			PetID:       luna.ID,
			MediaURL:    "https://images.unsplash.com/photo-1514888286974-6c03e2ca1dba?auto=format&fit=crop&w=800&q=80",
			Caption:     "Luna looking regal and waiting for her forever home at Happy Paws.",
			IsShowcased: true,
			CreatedAt:   time.Now().Add(-4 * time.Hour),
		},
		{
			ID:          uuid.New(),
			PetID:       luna.ID,
			MediaURL:    "https://images.unsplash.com/photo-1573865526739-10659fec78a5?auto=format&fit=crop&w=800&q=80",
			Caption:     "Purring session in the adoption lounge. She loves feather toys!",
			IsShowcased: true,
			CreatedAt:   time.Now().Add(-14 * time.Hour),
		},

		// Milo's Posts
		{
			ID:          uuid.New(),
			PetID:       milo.ID,
			MediaURL:    "https://images.unsplash.com/photo-1583511655857-d19b40a7a54e?auto=format&fit=crop&w=800&q=80",
			Caption:     "Mid-sprint action shot! Do not let the short legs fool you.",
			IsShowcased: true,
			CreatedAt:   time.Now().Add(-8 * time.Hour),
		},

		// Cleo's Posts
		{
			ID:          uuid.New(),
			PetID:       cleo.ID,
			MediaURL:    "https://images.unsplash.com/photo-1518791841217-8f162f1e1131?auto=format&fit=crop&w=800&q=80",
			Caption:     "Observing my kingdom from atop the bookshelf.",
			IsShowcased: true,
			CreatedAt:   time.Now().Add(-12 * time.Hour),
		},

		// Rocky's Posts
		{
			ID:          uuid.New(),
			PetID:       rocky.ID,
			MediaURL:    "https://images.unsplash.com/photo-1583337130417-3346a1be7dee?auto=format&fit=crop&w=800&q=80",
			Caption:     "New fall sweater arrived! How does my profile look?",
			IsShowcased: true,
			CreatedAt:   time.Now().Add(-15 * time.Hour),
		},

		// Oliver's Posts
		{
			ID:          uuid.New(),
			PetID:       oliver.ID,
			MediaURL:    "https://images.unsplash.com/photo-1548546738-8502ce9a3e0c?auto=format&fit=crop&w=800&q=80",
			Caption:     "Little Oliver is gaining weight and thriving after rescue! Ready for adoption soon.",
			IsShowcased: true,
			CreatedAt:   time.Now().Add(-1 * time.Hour),
		},
	}

	for i := range posts {
		database.DB.Create(&posts[i])
	}

	// 5. Seed Boops and Woofs
	log.Println("Seeding Boops and Woofs...")
	boops := []post.Boop{
		{ID: uuid.New(), AccountID: parentElena.ID, PostID: posts[0].ID, CreatedAt: time.Now()},
		{ID: uuid.New(), AccountID: parentMarcus.ID, PostID: posts[0].ID, CreatedAt: time.Now()},
		{ID: uuid.New(), AccountID: shelterHappyPaws.ID, PostID: posts[0].ID, CreatedAt: time.Now()},
		{ID: uuid.New(), AccountID: parentJohn.ID, PostID: posts[9].ID, CreatedAt: time.Now()},
		{ID: uuid.New(), AccountID: parentElena.ID, PostID: posts[9].ID, CreatedAt: time.Now()},
	}
	for i := range boops {
		database.DB.Create(&boops[i])
	}

	woofs := []post.Woof{
		{
			ID:        uuid.New(),
			AccountID: parentElena.ID,
			PostID:    posts[0].ID,
			Content:   "Buddy has the sweetest smile! Booped immediately.",
			CreatedAt: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:        uuid.New(),
			AccountID: parentMarcus.ID,
			PostID:    posts[0].ID,
			Content:   "Rocky wants to know if Buddy shares his sticks!",
			CreatedAt: time.Now().Add(-30 * time.Minute),
		},
		{
			ID:        uuid.New(),
			AccountID: parentJohn.ID,
			PostID:    posts[9].ID,
			Content:   "Luna is gorgeous! Hope she finds the loving family she deserves so much.",
			CreatedAt: time.Now().Add(-2 * time.Hour),
		},
	}
	for i := range woofs {
		database.DB.Create(&woofs[i])
	}

	// 6. Create Community Channels (Town Square)
	log.Println("Seeding Town Square Channels...")
	chanSOS := community.Channel{
		ID:           uuid.New(),
		AccountID:    shelterCityStrays.ID,
		Name:         "Tri-State SOS & Urgent Rescues",
		Description:  "Immediate rescue calls, lost pet alerts, and critical field response coordination.",
		CoverageArea: "New York, New Jersey, Connecticut",
		IsActive:     true,
	}

	chanAdoption := community.Channel{
		ID:           uuid.New(),
		AccountID:    shelterHappyPaws.ID,
		Name:         "Happy Paws Adoption Hub",
		Description:  "Meet adoptable cats, dogs, and small animals ready for their forever homes.",
		CoverageArea: "Greater NYC & Tri-State Area",
		IsActive:     true,
	}

	chanSanctuary := community.Channel{
		ID:           uuid.New(),
		AccountID:    shelterSecondChance.ID,
		Name:         "Sanctuary Care & Medical Drives",
		Description:  "Fundraisers, emergency veterinary care sponsorships, and winter shelter equipment.",
		CoverageArea: "Regional & Nationwide",
		IsActive:     true,
	}

	channels := []community.Channel{chanSOS, chanAdoption, chanSanctuary}
	for i := range channels {
		database.DB.Create(&channels[i])
	}

	// 7. Subscriptions
	log.Println("Seeding Subscriptions...")
	subs := []community.Subscription{
		{AccountID: parentJohn.ID, ChannelID: chanSOS.ID, CreatedAt: time.Now()},
		{AccountID: parentJohn.ID, ChannelID: chanAdoption.ID, CreatedAt: time.Now()},
		{AccountID: parentElena.ID, ChannelID: chanSOS.ID, CreatedAt: time.Now()},
		{AccountID: parentElena.ID, ChannelID: chanSanctuary.ID, CreatedAt: time.Now()},
		{AccountID: parentMarcus.ID, ChannelID: chanSOS.ID, CreatedAt: time.Now()},
	}
	for i := range subs {
		database.DB.Create(&subs[i])
	}

	// 8. Channel Posts (Town Square Alerts)
	log.Println("Seeding Town Square Alerts...")
	cPosts := []community.ChannelPost{
		{
			ID:        uuid.New(),
			ChannelID: chanSOS.ID,
			AuthorID:  shelterCityStrays.ID,
			Type:      "sos",
			Status:    "urgent",
			Title:     "Injured Kitten found near 5th Ave & 23rd St",
			Content:   "Emergency rescue needed! Tiny kitten with injured left leg discovered under scaffolding. We have a rescue carrier on site but need volunteer transport to Broadway Vet Clinic immediately! Emergency Dispatch Hotline: 555-0199.",
			MediaURL:  "https://images.unsplash.com/photo-1548546738-8502ce9a3e0c?auto=format&fit=crop&w=800&q=80",
			CreatedAt: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:        uuid.New(),
			ChannelID: chanSOS.ID,
			AuthorID:  shelterCityStrays.ID,
			Type:      "sos",
			Status:    "in_progress",
			Title:     "Lost Golden Retriever near Prospect Park Dog Beach",
			Content:   "Male Golden Retriever wearing red collar named Charlie ran off near 9th St entrance. Microchipped. Very friendly, responds to treats. Call 555-0142 if spotted.",
			MediaURL:  "https://images.unsplash.com/photo-1552053831-71594a27632d?auto=format&fit=crop&w=800&q=80",
			CreatedAt: time.Now().Add(-3 * time.Hour),
		},
		{
			ID:        uuid.New(),
			ChannelID: chanAdoption.ID,
			AuthorID:  shelterHappyPaws.ID,
			Type:      "adoption",
			Status:    "",
			Title:     "Meet Luna: Sweet 2-Year-Old Siamese Mix Seeking Quiet Home",
			Content:   "Luna is fully vaccinated, spayed, microchipped, and negative for FIV/FeLV. She loves chin scratches, feather wand toys, and curling up in warm sun puddles. Great with gentle adults and older pets!",
			MediaURL:  "https://images.unsplash.com/photo-1514888286974-6c03e2ca1dba?auto=format&fit=crop&w=800&q=80",
			CreatedAt: time.Now().Add(-12 * time.Hour),
		},
		{
			ID:        uuid.New(),
			ChannelID: chanAdoption.ID,
			AuthorID:  shelterHappyPaws.ID,
			Type:      "adoption",
			Status:    "",
			Title:     "Barnaby the Beagle Mix is Ready to Join Your Pack!",
			Content:   "Barnaby is a 1-year-old bundle of joy with soulful brown eyes and floppy ears. He loves group walks, agility obstacle runs, and cozy naps on the couch. Adoption fee includes full medical package.",
			MediaURL:  "https://images.unsplash.com/photo-1537151625747-768eb6cf92b2?auto=format&fit=crop&w=800&q=80",
			CreatedAt: time.Now().Add(-20 * time.Hour),
		},
		{
			ID:        uuid.New(),
			ChannelID: chanSanctuary.ID,
			AuthorID:  shelterSecondChance.ID,
			Type:      "donation",
			Status:    "",
			Title:     "Winter Blanket & Insulated Bedding Drive ($1,200 Goal)",
			Content:   "Winter temperatures are dropping below freezing. We are outfitting 40 sanctuary enclosures with heavy-duty thermal bedding and heated bowls. $25 covers one insulated set!",
			MediaURL:  "",
			CreatedAt: time.Now().Add(-24 * time.Hour),
		},
	}

	for i := range cPosts {
		database.DB.Create(&cPosts[i])
	}

	// 8. Seed Channel Admins (Designated Admins for Telegram-style Channels)
	log.Println("Seeding Designated Channel Admins...")
	admins := []community.ChannelAdmin{
		// Happy Paws Rescue Hub Admins
		{ChannelID: chanAdoption.ID, AccountID: shelterHappyPaws.ID, Role: "owner"},
		{ChannelID: chanAdoption.ID, AccountID: parentElena.ID, Role: "adoption_lead"},
		{ChannelID: chanAdoption.ID, AccountID: parentMarcus.ID, Role: "veterinary_admin"},

		// Tri-State SOS Admins
		{ChannelID: chanSOS.ID, AccountID: shelterCityStrays.ID, Role: "owner"},
		{ChannelID: chanSOS.ID, AccountID: parentJohn.ID, Role: "dispatch_volunteer"},

		// Second Chance Sanctuary Admins
		{ChannelID: chanSanctuary.ID, AccountID: shelterSecondChance.ID, Role: "owner"},
		{ChannelID: chanSanctuary.ID, AccountID: parentMarcus.ID, Role: "medical_coordinator"},
	}
	for i := range admins {
		database.DB.Create(&admins[i])
	}

	// 9. Seed Townhall Q&A Questions and Answers
	log.Println("Seeding Townhall Questions & Answers...")
	qEmergency := community.TownhallQuestion{
		ID:           uuid.New(),
		AuthorID:     parentJohn.ID,
		AuthorName:   "John Doe",
		AuthorType:   "parent",
		Title:        "URGENT: Dog ate 2 squares of 70% dark chocolate 20 mins ago. What should I do right now?",
		Content:      "My 25lb Golden retriever mix got into a dark chocolate bar on the coffee table. He seems fine right now but I know dark chocolate is toxic. Should I induce vomiting or rush straight to the animal hospital?",
		Category:     "urgent_health",
		IsUrgent:     true,
		HelpfulVotes: 48,
		CreatedAt:    time.Now().Add(-45 * time.Minute),
	}
	database.DB.Create(&qEmergency)

	aEmergency1 := community.TownhallAnswer{
		ID:           uuid.New(),
		QuestionID:   qEmergency.ID,
		AuthorID:     shelterHappyPaws.ID,
		AuthorName:   "Dr. Marcus Vance (Happy Paws Vet Lead)",
		AuthorType:   "shelter",
		Content:      "Do NOT induce vomiting with hydrogen peroxide at home without direct clinical supervision - it can cause severe gastric ulceration. For a 25lb dog, 2 squares (approx 20-30g) of 70% cocoa contains dangerous levels of theobromine. Please call ASPCA Animal Poison Control at (888) 426-4435 immediately or go directly to the nearest 24/7 veterinary emergency hospital. Bring the chocolate wrapper with you so the clinician can calculate exact mg/kg theobromine toxicity.",
		IsVerified:   true,
		HelpfulVotes: 64,
		CreatedAt:    time.Now().Add(-35 * time.Minute),
	}
	database.DB.Create(&aEmergency1)

	qDiet := community.TownhallQuestion{
		ID:           uuid.New(),
		AuthorID:     parentElena.ID,
		AuthorName:   "Elena Rostova",
		AuthorType:   "parent",
		Title:        "Best fiber and gut-friendly food recommendations for senior dogs with sensitive stomach?",
		Content:      "Looking for suggestions on balanced diets or digestive enzyme supplements for my 9-year-old pup who experiences occasional bouts of colitis after regular kibble.",
		Category:     "diet",
		IsUrgent:     false,
		HelpfulVotes: 19,
		CreatedAt:    time.Now().Add(-18 * time.Hour),
	}
	database.DB.Create(&qDiet)

	aDiet1 := community.TownhallAnswer{
		ID:           uuid.New(),
		QuestionID:   qDiet.ID,
		AuthorID:     parentMarcus.ID,
		AuthorName:   "Marcus Vance",
		AuthorType:   "parent",
		Content:      "Pure plain canned pumpkin puree (1-2 tablespoons per meal, make sure it has NO added spices or xylitol) worked wonders for our senior Frenchie. Also look into hydrolyzed protein veterinary formulas or adding a vet-grade probiotic like FortiFlora.",
		IsVerified:   false,
		HelpfulVotes: 23,
		CreatedAt:    time.Now().Add(-14 * time.Hour),
	}
	database.DB.Create(&aDiet1)

	qBehavior := community.TownhallQuestion{
		ID:           uuid.New(),
		AuthorID:     shelterCityStrays.ID,
		AuthorName:   "City Stray Alliance",
		AuthorType:   "shelter",
		Title:        "How to safely introduce a newly adopted rescue cat to an existing adult dog?",
		Content:      "We frequently get this question from our adopters! What has been the most successful scent-swapping and gradual door-barrier routine that worked for your multi-pet household?",
		Category:     "behavior",
		IsUrgent:     false,
		HelpfulVotes: 32,
		CreatedAt:    time.Now().Add(-2 * 24 * time.Hour),
	}
	database.DB.Create(&qBehavior)

	log.Println("Pawwp database successfully seeded with 6 profiles, 8 pets, 9-post showcases, channel admins, Town Square alerts, and Townhall Q&As!")
}
