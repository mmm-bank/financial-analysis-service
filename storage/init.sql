CREATE TABLE transactions (
    transaction_id UUID NOT NULL,
    user_id UUID NOT NULL,
    account_id UUID NOT NULL,
    partner_account_id UUID NOT NULL,
    card_id UUID,
    amount BIGINT NOT NULL,
    transaction_type VARCHAR(16) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),

    UNIQUE (transaction_id, transaction_type)
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id, created_at DESC);

CREATE OR REPLACE PROCEDURE add_transfer(
    id UUID,
    sender_id UUID,
    sender_account_id UUID,
    receiver_id UUID,
    receiver_account_id UUID,
    transfer_amount BIGINT
)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO transactions (transaction_id,
                              user_id,
                              account_id,
                              partner_account_id,
                              amount,
                              transaction_type)
    VALUES (id,
            sender_id,
            sender_account_id,
            receiver_account_id,
            transfer_amount,
            'send')
    ON CONFLICT (transaction_id, transaction_type) DO NOTHING;

    INSERT INTO  transactions (transaction_id,
                               user_id,
                               account_id,
                               partner_account_id,
                               amount,
                               transaction_type)
    VALUES (id,
            receiver_id,
            receiver_account_id,
            sender_account_id,
            transfer_amount,
            'receive')
    ON CONFLICT (transaction_id, transaction_type) DO NOTHING;

END;
$$