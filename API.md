# 🛰 Pawwp API Documentation

All protected routes require a `Bearer <access_token>` in the `Authorization` header.

## 🔐 Authentication & Accounts

| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| `POST` | `/accounts/register` | Create a new account | Public |
| `POST` | `/accounts/login` | Login and receive tokens | Public |
| `POST` | `/accounts/refresh` | Exchange refresh token for new tokens | Public |
| `GET` | `/accounts/me` | Get current user profile | Private |
| `PUT` | `/accounts/me` | Update current user profile (Bio) | Private |
| `DELETE` | `/accounts/me` | Delete account and all data | Private |

## 🐾 Pets

| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| `POST` | `/pets` | Register a new pet to your account | Private |
| `GET` | `/pets/{id}` | Get details of a specific pet | Private |
| `PUT` | `/pets/{id}` | Update pet details (Name, Bio, Status) | Private (Owner) |
| `DELETE` | `/pets/{id}` | Remove a pet | Private (Owner) |
| `GET` | `/accounts/{account_id}/pets` | List all pets owned by an account | Private |

## 📸 Posts & Feed (Phase 1 - In Progress)

| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| `POST` | `/posts` | Create a new post for a pet | Private |
| `GET` | `/feed` | Fetch the global 48-hour feed | Private |
| `GET` | `/pets/{pet_id}/showcase` | Fetch 9 showcased posts for a profile | Private |
| `PUT` | `/posts/{post_id}/showcase` | Toggle "Showcase" status (Max 9) | Private (Owner) |
| `GET` | `/pets/{pet_id}/archive` | Fetch all historical posts for a pet | Private (Owner) |
| `POST` | `/posts/{post_id}/boops` | Toggle a Boop (Like) | Private |
| `POST` | `/posts/{post_id}/woofs` | Add a Woof (Comment) | Private |

## 🚨 SOS & Announcements (Upcoming)

| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| `GET` | `/sos` | Fetch active SOS alerts | Private |
| `POST` | `/sos` | Post a new emergency alert | Private |
| `GET` | `/announcements` | Fetch adoption/donation feed | Private |
| `POST` | `/announcements` | Post a new adoption or donation call | Private |

## ✉️ Messaging (Upcoming)

| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| `GET` | `/dms` | List recent conversations | Private |
| `POST` | `/dms/{account_id}` | Send a message | Private |
| `GET` | `/dms/{account_id}` | Fetch message history | Private |
