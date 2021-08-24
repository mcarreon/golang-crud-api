create database if not exists postgres_db

create table if not exists postgres_db.books (
  id UUID primary key default gen_random_uuid(),
  title text not null,
  author text not null,
  published_date date ,
  publisher text not null,
  rating int not NULL,
  status text not null
)