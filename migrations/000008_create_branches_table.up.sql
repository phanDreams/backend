CREATE TABLE IF NOT EXISTS branches (
    id         INTEGER PRIMARY KEY,
    address_id INTEGER,
    is_hq      BOOLEAN,
    phone      VARCHAR(50),
    
    CONSTRAINT fk_branches_address
        FOREIGN KEY (address_id) REFERENCES addresses(id)
);
