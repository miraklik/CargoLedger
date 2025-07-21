CREATE TYPE statusEn AS ENUM ('created', 'inTransit', 'delivered', 'cancelled');

CREATE TABLE cargos (
    ID BIGSERIAL PRIMARY KEY,
    Sender VARCHAR(255) NOT NULL,
    Carrier VARCHAR(255) NOT NULL,
    Receiver VARCHAR(255) NOT NULL,
    Ipfs_Hash TEXT NOT NULL,
    Status statusEn
);

CREATE INDEX idx_cargos_id ON cargos(ID);