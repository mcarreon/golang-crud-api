CREATE TABLE IF NOT EXISTS books (
  id UUID primary key default gen_random_uuid(),
  title text not null,
  author text not null,
  published_date date,
  publisher text not null,
  rating int not NULL,
  status text not null
)

