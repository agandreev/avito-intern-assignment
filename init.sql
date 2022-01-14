CREATE TABLE users
(
    id      SERIAL PRIMARY KEY,
    user_id INT,
    amount  NUMERIC
);

CREATE TABLE operations
(
    id           SERIAL PRIMARY KEY,
    initiator_id INT,
    type         VARCHAR(20),
    amount       NUMERIC,
    time         TIMESTAMP,
    receiver_id  INT,
    FOREIGN KEY (initiator_id) REFERENCES users(id),
    FOREIGN KEY (receiver_id) REFERENCES users(id)
);