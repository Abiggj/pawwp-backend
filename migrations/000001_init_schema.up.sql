CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Accounts
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(20) NOT NULL CHECK (type IN ('individual', 'shelter', 'zoo', 'breeder')),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    bio TEXT,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_accounts_type ON accounts(type);

-- Pets
CREATE TABLE pets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    username VARCHAR(30) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    species VARCHAR(50) NOT NULL,
    breed VARCHAR(100),
    birth_date DATE,
    bio TEXT,
    avatar_url TEXT,
    adoption_status VARCHAR(20)
        CHECK (adoption_status IN ('not_applicable', 'adoptable', 'adopted', 'rescued')),
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_pets_account_id ON pets(account_id);
CREATE INDEX idx_pets_species ON pets(species);

-- Posts
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID REFERENCES pets(id) ON DELETE CASCADE,
    account_id UUID REFERENCES accounts(id) ON DELETE CASCADE,
    caption TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    CHECK (
        (pet_id IS NOT NULL AND account_id IS NULL)
        OR
        (pet_id IS NULL AND account_id IS NOT NULL)
    )
);

CREATE INDEX idx_posts_pet_id ON posts(pet_id);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);

-- Media
CREATE TABLE media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    media_url TEXT NOT NULL,
    media_type VARCHAR(20) CHECK (media_type IN ('image', 'video')),
    position INT DEFAULT 0
);

-- Follows
CREATE TABLE follows (
    follower_pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    following_pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    PRIMARY KEY (follower_pet_id, following_pet_id),
    CHECK (follower_pet_id <> following_pet_id)
);

CREATE INDEX idx_follows_following ON follows(following_pet_id);

-- Likes
CREATE TABLE likes (
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    PRIMARY KEY (pet_id, post_id)
);

-- Comments
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    parent_comment_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_comments_post_id ON comments(post_id);

-- Notifications
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipient_pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    type VARCHAR(30) NOT NULL,
    actor_pet_id UUID REFERENCES pets(id) ON DELETE CASCADE,
    post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
