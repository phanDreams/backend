CREATE TABLE IF NOT EXISTS service_animals (
    service_id INTEGER,
    animal_id  INTEGER,
    PRIMARY KEY (service_id, animal_id),

    CONSTRAINT fk_sa_service
        FOREIGN KEY (service_id) REFERENCES services(id),
    CONSTRAINT fk_sa_animal
        FOREIGN KEY (animal_id) REFERENCES animals(id)
);