CREATE TABLE IF NOT EXISTS users_accounts(
    user_id INT NOT NULL,
    account_id INT NOT NULL,
    PRIMARY KEY (user_id, account_id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_account_id FOREIGN KEY (account_id) REFERENCES accounts (id)
);