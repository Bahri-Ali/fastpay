DROP TYPE IF EXISTS user_role CASCADE;
CREATE TYPE user_role AS ENUM('normal', 'merchant', 'child');

DROP TYPE IF EXISTS transaction_type CASCADE;
CREATE TYPE transaction_type AS ENUM('deposit', 'withdrawal', 'transfer', 'payment');

DROP TYPE IF EXISTS transaction_status CASCADE;
CREATE TYPE transaction_status AS ENUM('pending', 'completed', 'failed');

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  phone_number TEXT UNIQUE,
  email TEXT UNIQUE,
  password_hash TEXT,
  full_name TEXT,
  role user_role,
  parent_id UUID REFERENCES users (id),
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wallets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  user_id UUID UNIQUE REFERENCES users (id),
  balance NUMERIC(19, 4) DEFAULT 0.0,
  currency TEXT DEFAULT 'DZD',
  is_frozen BOOLEAN DEFAULT false,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  wallet_id_from UUID REFERENCES wallets (id),
  wallet_id_to UUID REFERENCES wallets (id),
  amount NUMERIC(19, 4),
  type transaction_type,
  status transaction_status,
  reference_code TEXT UNIQUE,
  idempotency_key TEXT UNIQUE,
  description TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS merchants (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  user_id UUID UNIQUE REFERENCES users (id),
  business_name TEXT,
  business_address TEXT,
  verification_status TEXT
);

CREATE TABLE IF NOT EXISTS cards (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  user_id UUID REFERENCES users (id),
  card_pan_hash TEXT,
  last_four_digits TEXT,
  expiry_date DATE,
  cvv_hash TEXT,
  card_limit NUMERIC(19, 4),
  is_active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS contacts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  user_id UUID REFERENCES users (id),
  contact_user_id UUID REFERENCES users (id),
  alias TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS service_providers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  name TEXT,
  category TEXT,
  logo_url TEXT
);

CREATE TABLE IF NOT EXISTS linked_bills (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  user_id UUID REFERENCES users (id),
  provider_id UUID REFERENCES service_providers (id),
  account_reference TEXT,
  auto_pay BOOLEAN DEFAULT false
);

CREATE TABLE IF NOT EXISTS spending_limits (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  user_id UUID REFERENCES users (id),
  daily_limit NUMERIC(19, 4),
  monthly_limit NUMERIC(19, 4)
);

CREATE TABLE IF NOT EXISTS notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
  user_id UUID REFERENCES users (id),
  title TEXT,
  message TEXT,
  is_read BOOLEAN DEFAULT false,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_phone_number ON users USING btree (phone_number);
CREATE INDEX IF NOT EXISTS idx_transactions_reference_code ON transactions USING btree (reference_code);