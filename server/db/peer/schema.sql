CREATE TABLE IF NOT EXISTS peers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    peerID TEXT UNIQUE NOT NULL,
    multiaddrs TEXT NOT NULL,
    pubKey TEXT NOT NULL,
    timestamp BIGINT NOT NULL,
    signature TEXT NOT NULL,
    status INT NOT NULL
);

CREATE INDEX idx_peers_pID ON peers(peerID);