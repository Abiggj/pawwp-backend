# 🛰 Pawwp API Documentation

All protected routes require a `Bearer <access_token>` in the `Authorization` header.

## 🔐 Authentication & Accounts

| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| `POST` | `/accounts/register` | Create a new account (`parent` or `shelter`) | Public |
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

## 📸 Personal Media (FYP & Showcase)

| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| `POST` | `/posts` | Create a 48h post for a pet | Private |
| `GET` | `/feed` | Fetch the global 48-hour FYP feed | Private |
| `GET` | `/pets/{pet_id}/showcase` | Fetch 9 showcased posts for a profile | Private |
| `PUT` | `/posts/{post_id}/showcase` | Toggle "Showcase" status (Max 9) | Private (Owner) |
| `GET` | `/pets/{pet_id}/archive` | Fetch all historical posts for a pet | Private (Owner) |
| `POST` | `/posts/{post_id}/boops` | Toggle a Boop (Like) | Private |
| `POST` | `/posts/{post_id}/woofs` | Add a Woof (Comment) | Private |

## 🏘 Community Channels (SOS, Adoption, Donation)

| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| `POST` | `/community/channels` | Create a new Channel (Shelters only) | Private |
| `GET` | `/community/channels` | Browse all active channels | Private |
| `GET` | `/community/channels/{id}` | Get channel details | Private |
| `PUT` | `/community/channels/{id}` | Update channel (Owner only) | Private |
| `POST` | `/community/channels/{id}/posts` | Post to a channel (`type`: `sos`, `adoption`, `donation`) | Private |
| `GET` | `/community/channels/{id}/posts` | View all posts in a channel | Private |
| `PUT` | `/community/posts/{post_id}/status` | Update SOS status (urgent, resolved, etc.) | Private |
| `POST` | `/community/channels/{id}/subscribe` | Subscribe to a channel | Private |
| `DELETE` | `/community/channels/{id}/subscribe` | Unsubscribe from a channel | Private |
| `GET` | `/accounts/me/subscriptions` | List my subscribed channels | Private |
| `GET` | `/community/feed` | Feed of posts from subscribed channels | Private |

## ✉️ Messaging (Upcoming)

| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| `GET` | `/dms` | List recent conversations | Private |
| `POST` | `/dms/{account_id}` | Send a message | Private |
| `GET` | `/dms/{account_id}` | Fetch message history | Private |
