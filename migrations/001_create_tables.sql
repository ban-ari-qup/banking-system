CREATE TABLE IF NOT EXISTS accounts (
    id TEXT PRIMARY KEY,
    password TEXT NOT NULL,
    cvc2 TEXT NOT NULL,
    balance DECIMAL(15,2) NOT NULL DEFAULT 0,
    name TEXT NOT NULL,
    phone TEXT UNIQUE NOT NULL,
    age INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expired_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    type TEXT NOT NULL,
    from_account TEXT,
    to_account TEXT,
    amount DECIMAL(15,2) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    status TEXT NOT NULL,
    account_id TEXT NOT NULL,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_activity TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    user_agent TEXT NOT NULL,
    ip TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES accounts(id) ON DELETE CASCADE
);