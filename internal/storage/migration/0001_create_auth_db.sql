CREATE TABLE user
(
    id           BIGSERIAL PRIMARY KEY NOT NULL,
    name         VARCHAR(64) NOT NULL,
    email        VARCHAR(64) NOT NULL,
    password     VARCHAR(64),
    registeredAt TIMESTAMP DEFAULT now(),
    lastVisitAt  TIMESTAMP
);
