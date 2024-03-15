CREATE TABLE "book" (
    id SERIAL PRIMARY KEY,
    isbn VARCHAR(100) NOT NULL,
    "name" VARCHAR(200) NOT NULL,
    author VARCHAR(200) NOT NULL,
    publish_year SMALLINT NOT NULL
);

INSERT INTO book(isbn, "name", author, publish_year) 
VALUES 
    ('978-93-5019-561-1', 'Junior Level Books Introduction to Computer', 'Amit Garg', 2011),
    ('978-93-8067-432-2', 'Client Server Computing', 'Lalit Kumar', 2012),
    ('978-93-5163-389-1 ', 'Data Structure Using C', 'Sharad Kumar Verma', 2015),
    ('978-93-8067-432-2', ' Client Server Computing', 'Sharad Kumar Verma', 2012);
