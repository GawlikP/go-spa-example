package query

const ClearTestsDatabaseQuerry = `
  DELETE FROM posts;
  ALTER SEQUENCE posts_id_seq RESTART WITH 1;
  DELETE FROM users;
  ALTER SEQUENCE users_id_seq RESTART WITH 1;
`
