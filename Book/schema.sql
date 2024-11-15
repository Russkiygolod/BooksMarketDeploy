
DROP TABLE IF EXISTS books, authors, authors_books;

CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    price INTEGER NOT NULL
);

CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL
);

CREATE TABLE authors_books (
    books_id INTEGER REFERENCES books(id),
    authors_id INTEGER REFERENCES authors(id)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL 
)
