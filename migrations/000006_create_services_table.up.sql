CREATE TABLE IF NOT EXISTS services (
    id                INTEGER PRIMARY KEY,
    name              VARCHAR(100) NOT NULL,
    category_id       INTEGER,
    animalcategory_id INTEGER,
    
    CONSTRAINT fk_services_category
        FOREIGN KEY (category_id) REFERENCES categories(id),
    CONSTRAINT fk_services_animal_category
        FOREIGN KEY (animalcategory_id) REFERENCES animal_categories(id)
);