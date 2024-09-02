# resets db, then migrates where up again

cd sql/schema
goose postgres postgres://postgres:postgres@localhost:5432/blogator reset
goose postgres postgres://postgres:postgres@localhost:5432/blogator up