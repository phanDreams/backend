-- Medical services
INSERT INTO services (name, category_id, duration) VALUES
('Консультація ветеринара', 1, 30),
('Вакцинація', 1, 15),
('Стерилізація/кастрація', 1, 120),
('Чистка зубів', 1, 60),
('Діагностика', 1, 45);

-- Grooming services
INSERT INTO services (name, category_id, duration) VALUES
('Стрижка', 2, 60),
('Купання', 2, 30),
('Чистка вух', 2, 15),
('Стрижка кігтів', 2, 15),
('Вичісування', 2, 30);

-- Training services
INSERT INTO services (name, category_id, duration) VALUES
('Базове навчання', 3, 60),
('Корекція поведінки', 3, 60),
('Навчання трюкам', 3, 60),
('Групові заняття', 3, 90); 