DROP TABLE IF EXISTS book_author;

CREATE TABLE book_author (
    book_id INT,
    author_id INT,
    FOREIGN KEY (book_id) REFERENCES book(id),
    FOREIGN KEY (author_id) REFERENCES author(id)
);

INSERT INTO book_author (book_id, author_id)
SELECT b.id, a.id
FROM book b
JOIN author a ON b.author_id = a.id;

ALTER TABLE book DROP COLUMN author_id;