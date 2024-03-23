package query

const AddUsersQuerry = `
  INSERT INTO users
  (email, password, nickname)
  VALUES ($1, $2, $3)
  RETURNING
  id, email, password, nickname, created_at, updated_at
`
