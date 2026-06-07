-- TODO: разделить запросы по разным SQL-скриптам
DROP TABLE IF EXISTS tasks;

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO tasks (title, description, completed)
VALUES (
        'Пошкодничать',
        'Подёргать микукшу за косички',
        TRUE
    ),
    (
        'Написать REST API',
        'Самому написать, ага',
        FALSE
    ),
    (
        'Релиз аппки',
        'Деплойку сделать',
        FALSE
    );