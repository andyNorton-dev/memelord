DROP TABLE IF EXISTS user_workers;
DROP TABLE IF EXISTS workers_upgrade;
DROP TABLE IF EXISTS workers;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    tg_id BIGINT NOT NULL UNIQUE,
    username VARCHAR(255),
    balance INT DEFAULT 0,
    level INT DEFAULT 1,
    energy INT DEFAULT 10,
    max_energy INT DEFAULT 10,
    profit_per_hour INT DEFAULT 0,
    head TEXT,
    body TEXT,
    legs TEXT,
    foot TEXT,
    profit_for_tap INT DEFAULT 1,
    last_restoration TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_profit_per_hour TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE workers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(45) NOT NULL,
    description TEXT,
    url_image TEXT
);

CREATE TABLE workers_upgrade (
    id SERIAL PRIMARY KEY,
    id_worker INT NOT NULL,
    level INT NOT NULL,
    cost INT NOT NULL,
    profit INT NOT NULL,
    FOREIGN KEY (id_worker) REFERENCES workers(id)
);

CREATE TABLE user_workers (
    id SERIAL PRIMARY KEY,
    id_worker INT NOT NULL,
    id_upgrade INT NOT NULL,
    id_user INT NOT NULL,
    FOREIGN KEY (id_worker) REFERENCES workers(id),
    FOREIGN KEY (id_upgrade) REFERENCES workers_upgrade(id)
); 