DROP TABLE IF EXISTS author;

CREATE TABLE "author" (
    id SERIAL PRIMARY KEY,
    "name" VARCHAR(200) NOT NULL,
    birth_day DATE
);

INSERT INTO author("name", birth_day)
SELECT "author", null from book;

INSERT INTO author("name", birth_day) 
VALUES 
    ('Amit Garg', '1978-03-18'),
    ('Lalit Kumar', '1970-01-01'),
    ('Sharad Kumar Verma', '1987-06-24'),
    ('James Clear', '1986-01-22');

ALTER TABLE book
    ADD COLUMN author_id INT,
    ADD CONSTRAINT fk_author FOREIGN KEY (author_id) REFERENCES author(id);

UPDATE book
    SET author_id = author.id
    FROM author
    WHERE book.author = author.name;

ALTER TABLE book
    DROP COLUMN author;