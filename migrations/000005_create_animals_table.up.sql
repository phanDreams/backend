CREATE TABLE IF NOT EXISTS animals (
    id                INTEGER PRIMARY KEY,
    animalcategory_id INTEGER,
    size              animal_size,     -- e.g., 'small', 'medium', 'large', 'giant'
    breed             VARCHAR(100),

    CONSTRAINT fk_animals_animal_category
        FOREIGN KEY (animalcategory_id) REFERENCES animal_categories(id)
);