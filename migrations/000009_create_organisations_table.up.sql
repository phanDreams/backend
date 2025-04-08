CREATE TABLE organisations (
    id            INTEGER PRIMARY KEY,
    name          VARCHAR(100),
    email         VARCHAR(200),
    password_hash VARCHAR(200),
    description   TEXT,
    avatar        TEXT,
    is_banned     BOOLEAN,
    is_deleted    BOOLEAN,
    image_id      TEXT[],
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ
);

CREATE TABLE organisation_services (
    organisation_id INTEGER,
    service_id      INTEGER,
    PRIMARY KEY (organisation_id, service_id),
    CONSTRAINT fk_org_services_org
        FOREIGN KEY (organisation_id) REFERENCES organisations(id),
    CONSTRAINT fk_org_services_service
        FOREIGN KEY (service_id) REFERENCES services(id)
);

CREATE TABLE organisation_branches (
    organisation_id INTEGER,
    branch_id       INTEGER,
    PRIMARY KEY (organisation_id, branch_id),
    CONSTRAINT fk_org_branches_org
        FOREIGN KEY (organisation_id) REFERENCES organisations(id),
    CONSTRAINT fk_org_branches_branch
        FOREIGN KEY (branch_id) REFERENCES branches(id)
);
