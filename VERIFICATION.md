# Pawwp Backend Verification Guide

This guide provides CURL commands to test the core features of the Pawwp backend.

**Note:** You will need to replace `ACCESS_TOKEN` and `REFRESH_TOKEN` with the values received from the login response.

## 1. Authentication & Accounts

### Register a Parent
```bash
curl -X POST http://localhost:8080/accounts/register \
-H "Content-Type: application/json" \
-d '{
  "type": "parent",
  "email": "test_parent@example.com",
  "password": "password123",
  "bio": "I love golden retrievers"
}'
```

### Register a Shelter
```bash
curl -X POST http://localhost:8080/accounts/register \
-H "Content-Type: application/json" \
-d '{
  "type": "shelter",
  "email": "test_shelter@example.com",
  "password": "password123",
  "bio": "Dedicated to animal welfare"
}'
```

### Login
```bash
curl -X POST http://localhost:8080/accounts/login \
-H "Content-Type: application/json" \
-d '{
  "email": "test_parent@example.com",
  "password": "password123"
}'
```

### Refresh Token
```bash
curl -X POST http://localhost:8080/accounts/refresh \
-H "Content-Type: application/json" \
-d '{
  "refresh_token": "YOUR_REFRESH_TOKEN"
}'
```

---

## 2. Pet Management

### Create a Pet
```bash
curl -X POST http://localhost:8080/pets \
-H "Authorization: Bearer ACCESS_TOKEN" \
-H "Content-Type: application/json" \
-d '{
  "username": "buddy_goldie",
  "name": "Buddy",
  "species": "Dog",
  "breed": "Golden Retriever",
  "bio": "Loves fetch!",
  "adoption_status": "not_applicable"
}'
```

---

## 3. Personal Media (FYP & Showcase)

### Create a Post
```bash
curl -X POST http://localhost:8080/posts \
-H "Authorization: Bearer ACCESS_TOKEN" \
-H "Content-Type: application/json" \
-d '{
  "pet_id": "YOUR_PET_ID",
  "media_url": "https://example.com/photo.jpg",
  "caption": "Sunny day!"
}'
```

### Get Global 48h Feed (FYP)
```bash
curl -X GET http://localhost:8080/feed \
-H "Authorization: Bearer ACCESS_TOKEN"
```

### Toggle Showcase (Pin to Profile)
```bash
curl -X PUT http://localhost:8080/posts/YOUR_POST_ID/showcase \
-H "Authorization: Bearer ACCESS_TOKEN"
```

---

## 4. Community Channels

### Create a Channel (Shelter only)
```bash
curl -X POST http://localhost:8080/community/channels \
-H "Authorization: Bearer SHELTER_ACCESS_TOKEN" \
-H "Content-Type: application/json" \
-d '{
  "name": "Central Rescue Hub",
  "description": "City-wide emergency alerts",
  "coverage_area": "New York"
}'
```

### Post SOS to Channel
```bash
curl -X POST http://localhost:8080/community/channels/YOUR_CHANNEL_ID/posts \
-H "Authorization: Bearer ACCESS_TOKEN" \
-H "Content-Type: application/json" \
-d '{
  "type": "sos",
  "title": "Injured Kitten",
  "content": "Emergency rescue needed at Central Park",
  "status": "urgent"
}'
```

### Subscribe to Channel
```bash
curl -X POST http://localhost:8080/community/channels/YOUR_CHANNEL_ID/subscribe \
-H "Authorization: Bearer ACCESS_TOKEN"
```

### Get Subscribed Feed
```bash
curl -X GET http://localhost:8080/community/feed \
-H "Authorization: Bearer ACCESS_TOKEN"
```
