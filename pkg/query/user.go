package query

const AddUser = `
  INSERT INTO users
  (email, password, nickname)
  VALUES ($1, $2, $3)
  RETURNING
  id, email, password, nickname, created_at, updated_at
`

const FindUserByEmailAndNickname = `
  SELECT *
  FROM users
  WHERE
  users.email = $1
  OR
  users.nickname = $2
  LIMIT 1
`

const FindUser = `
  SELECT *
  FROM users
  WHERE
  users.id = $1
`
