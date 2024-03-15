ALTER TABLE book
    ADD COLUMN author VARCHAR(200);

UPDATE book
    SET author = author.name
    FROM author
    WHERE book.author_id = author.id;

ALTER TABLE book
    DROP CONSTRAINT fk_author;
    
DELETE FROM author;
DROP TABLE author;

ALTER TABLE book
    DROP COLUMN author_id;