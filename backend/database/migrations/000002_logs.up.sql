CREATE TABLE logs(
    ID BIGSERIAL PRIMARY KEY,
    Event_type VARCHAR(50) NOT NULL,
    Cargo_id BIGINT,
    Tx_hash VARCHAR(255) NOT NULL,
    Block_num BIGINT NOT NULL,
    Event_Index INT,
    Sender_Address VARCHAR(50),
    Data JSONB,
);