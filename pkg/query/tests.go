package query

const ClearTestsDatabaseQuerry = `
  DELETE FROM users;
  ALTER SEQUENCE users_id_seq RESTART WITH 1;
`
