package query

const AddPost = `
  INSERT INTO posts
  (title, content, user_id)
  VALUES
  ($1, $2, $3)
  RETURNING id, title, content, user_id, created_at, updated_at, deleted
`

const CountPosts = `
  SELECT COUNT(*)
  FROM posts
`

const AllPostsWithUsers = `
  SELECT
    p.id, p.title, p.content, p.user_id, p.created_at, p.updated_at, p.deleted,
    u.id, u.email, u.nickname, u.created_at, u.updated_at
  FROM posts p
  JOIN users u ON p.user_id = u.id
  ORDER BY p.created_at DESC
  LIMIT $1 OFFSET $2
`

const UpdatePost = `
  UPDATE posts
  SET title = $1, content = $2, user_id = $3, updated_at = NOW()
  WHERE id = $4
  RETURNING id, title, content, user_id, created_at, updated_at, deleted
`

const FindPostByID = `
  SELECT
    p.id, p.title, p.content, p.user_id, p.created_at, p.updated_at, p.deleted,
    u.id, u.email, u.nickname, u.created_at, u.updated_at
  FROM posts p
  JOIN users u ON p.user_id = u.id
  WHERE p.id = $1
`

const DeletePost = `
  UPDATE posts
  SET deleted = true
  WHERE id = $1
  RETURNING id, title, content, user_id, created_at, updated_at, deleted
`
