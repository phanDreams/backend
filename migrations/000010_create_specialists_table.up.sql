CREATE TABLE IF NOT EXISTS specialists (
    id              INTEGER PRIMARY KEY,
    name            VARCHAR(100),
    family_name     VARCHAR(100),
    email           VARCHAR(200) UNIQUE,
    password_hash   VARCHAR(200),
    bio             TEXT,
    avatar          TEXT,
    address_id      INTEGER,
    organisation_id INTEGER,
    is_banned       BOOLEAN,
    is_deleted      BOOLEAN,
    image_id        TEXT[],
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ,
    description     TEXT,
    phone           VARCHAR(50),

    CONSTRAINT fk_specialists_address
        FOREIGN KEY (address_id) REFERENCES addresses(id),
    CONSTRAINT fk_specialists_organisation
        FOREIGN KEY (organisation_id) REFERENCES organisations(id)
);

CREATE TABLE IF NOT EXISTS specialist_services (
    specialist_id INTEGER,
    service_id    INTEGER,
    PRIMARY KEY (specialist_id, service_id),
    CONSTRAINT fk_specialist_services_specialist
        FOREIGN KEY (specialist_id) REFERENCES specialists(id),
    CONSTRAINT fk_specialist_services_service
        FOREIGN KEY (service_id) REFERENCES services(id)
);
