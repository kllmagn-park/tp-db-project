package queries

const (
	CreateThreadCommand     = "INSERT INTO threads (title, author, message, created, slug, forum) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;"
	GetThreadByIdCommand    = "SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE id = $1;"
	GetThreadBySlugCommand  = "SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE slug = $1;"
	UpdateThreadByIdCommand = "UPDATE threads SET (title, message) = ($1, $2) WHERE id = $3;"

	GetPostsOnThreadFlatCommand                    = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts WHERE thread = $1 AND id > $2 ORDER BY created, id LIMIT $3;"
	GetPostsOnThreadFlatDescCommand                = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts WHERE thread = $1 AND id < $2 ORDER BY created DESC, id DESC LIMIT $3;"
	GetPostsOnThreadTreeCommand                    = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts WHERE thread = $1 AND parent_path > (SELECT parent_path FROM posts WHERE id = $2) ORDER BY parent_path, id LIMIT $3;"
	GetPostsOnThreadTreeDescCommand                = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts WHERE thread = $1 AND parent_path < (SELECT parent_path FROM posts WHERE id = $2) ORDER BY parent_path DESC LIMIT $3;"
	GetPostsOnThreadParentTreeCommand              = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts WHERE parent_path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0 AND id > (SELECT parent_path[1] FROM posts WHERE id = $2) ORDER BY id LIMIT $3) ORDER BY parent_path, id;"
	GetPostsOnThreadParentTreeDescWithSinceCommand = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts WHERE parent_path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0 AND id < (SELECT parent_path[1] FROM posts WHERE id = $2) ORDER BY id DESC LIMIT $3) ORDER BY parent_path[1] DESC, parent_path, id;"

	GetPostsOnThreadFlatWithoutSinceCommand           = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts WHERE thread = $1 ORDER BY created, id LIMIT $2;"
	GetPostsOnThreadFlatDescWithoutSinceCommand       = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts WHERE thread = $1 ORDER BY created DESC, id DESC LIMIT $2;"
	GetPostsOnThreadTreeWithoutSinceCommand           = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts WHERE thread = $1 ORDER BY parent_path, id LIMIT $2;"
	GetPostsOnThreadTreeDescWithoutSinceCommand       = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts WHERE thread = $1 ORDER BY parent_path DESC LIMIT $2;"
	GetPostsOnThreadParentTreeWithoutSinceCommand     = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts WHERE parent_path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0 ORDER BY id LIMIT $2) ORDER BY parent_path, id;"
	GetPostsOnThreadParentTreeDescWithoutSinceCommand = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts WHERE parent_path[1] IN (SELECT id FROM posts WHERE thread = $1 AND parent = 0 ORDER BY id DESC LIMIT $2) ORDER BY parent_path[1] DESC, parent_path, id;"
	
	DeleteTablesCommand    = "TRUNCATE TABLE users, forums, threads, posts, forumToUsers, votes CASCADE;"
	GetCountRecordsCommand = "SELECT (SELECT count(*) FROM users), (SELECT count(*) FROM forums), (SELECT count(*) FROM threads), (SELECT count(*) FROM posts);"

	CreateForumCommand = "INSERT INTO forums (title, \"user\", slug) VALUES ($1, $2, $3);"
	GetForumCommand    = "SELECT title, \"user\", slug, posts, threads FROM forums WHERE slug = $1;"

	GetUsersOnForumCommand                 = "SELECT nickname, fullname, about, email FROM forumToUsers WHERE forum = $1 AND nickname > $2 ORDER BY nickname LIMIT $3;"
	GetUsersOnForumDescCommand             = "SELECT nickname, fullname, about, email FROM forumToUsers WHERE forum = $1 AND nickname < $2 ORDER BY nickname DESC LIMIT $3;"
	GetUsersOnForumWithoutSinceCommand     = "SELECT nickname, fullname, about, email FROM forumToUsers WHERE forum = $1 ORDER BY nickname LIMIT $2;"
	GetUsersOnForumWithoutSinceDescCommand = "SELECT nickname, fullname, about, email FROM forumToUsers WHERE forum = $1 ORDER BY nickname DESC LIMIT $2;"

	GetThreadsOnForumCommand                 = "SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE forum = $1 AND created >= $2 ORDER BY created LIMIT $3;"
	GetThreadsOnForumDescCommand             = "SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE forum = $1 AND created <= $2 ORDER BY created DESC LIMIT $3;"
	GetThreadsOnForumWithoutSinceCommand     = "SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE forum = $1 ORDER BY created LIMIT $2;"
	GetThreadsOnForumWithoutSinceDescCommand = "SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE forum = $1 ORDER BY created DESC LIMIT $2;"

	GetPostCommand       = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts WHERE id = $1;"
	GetPostAuthorCommand = "SELECT nickname, fullname, about, email FROM users WHERE nickname = $1;"
	GetPostForumCommand  = "SELECT title, \"user\", slug, posts, threads FROM forums WHERE slug = $1;"
	GetPostThreadCommand = "SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE id = $1;"
	UpdatePostCommand    = "UPDATE posts SET (message, isEdited) = ($1, true) WHERE id = $2;"

	CreateUserCommand               = "INSERT INTO users (nickname, fullname, about, email) VALUES ($1, $2, $3, $4);"
	UpdateUserCommand               = "UPDATE users SET (fullname, about, email) = ($1, $2, $3) WHERE nickname = $4;"
	GetUserByNicknameCommand        = "SELECT nickname, fullname, about, email FROM users WHERE nickname = $1;"
	GetUserByEmailCommand           = "SELECT nickname, fullname, about, email FROM users WHERE email = $1;"
	GetUserByNicknameOrEmailCommand = "SELECT nickname, fullname, about, email FROM users WHERE nickname = $1 OR email = $2;"

	GetVoteByNicknameAndThreadCommand = "SELECT nickname, thread, voice FROM votes WHERE nickname = $1 AND thread = $2;"
	CreateVoteCommand                 = "INSERT INTO votes (nickname, thread, voice) VALUES ($1, $2, $3);"
	UpdateVoteCommand                 = "UPDATE votes SET voice = $1 WHERE nickname = $2 AND thread = $3 AND voice != $1;"
)