-- AMI Schema for DoltDB

CREATE TABLE IF NOT EXISTS memories (
    id VARCHAR(36) PRIMARY KEY,
    content TEXT NOT NULL,
    category ENUM('core', 'semantic', 'working', 'episodic') DEFAULT 'episodic',
    priority FLOAT DEFAULT 0.5,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    access_count INT DEFAULT 0,
    source VARCHAR(255),
    tags JSON,
    embedding BLOB
);

CREATE TABLE IF NOT EXISTS memory_links (
    from_id VARCHAR(36),
    to_id VARCHAR(36),
    relation VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (from_id, to_id, relation),
    FOREIGN KEY (from_id) REFERENCES memories(id),
    FOREIGN KEY (to_id) REFERENCES memories(id)
);

CREATE INDEX idx_memories_category ON memories(category);
CREATE INDEX idx_memories_priority ON memories(priority DESC);
CREATE INDEX idx_memories_accessed ON memories(accessed_at DESC);
