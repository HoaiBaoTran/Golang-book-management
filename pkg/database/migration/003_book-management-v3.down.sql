ALTER TABLE book
ADD COLUMN author_id INT;

UPDATE book b
SET author_id = ba.author_id
FROM book_author ba
WHERE b.id = ba.book_id;

DROP TABLE IF EXISTS book_author CASCADE;