package postcomments

import "database/sql"

type PostComment struct {
	PostID    int `json:"post_id"`
	CommentID int `json:"comment_id"`
}

type Comment struct {
	ID       int
	Content  string
	RegDate  string
	UserID   int
	Username string
	Avatar   string
}

func CreatePostComment(db *sql.DB, postComment PostComment) error {
	query := `INSERT INTO post_comments (post_id, comment_id)
              VALUES ($1, $2)`
	_, err := db.Exec(query, postComment.PostID, postComment.CommentID)
	return err
}

func GetPostComment(db *sql.DB, postID, commentID int) (*PostComment, error) {
	var postComment PostComment
	query := `SELECT post_id, comment_id FROM post_comments WHERE post_id=$1 AND comment_id=$2`
	err := db.QueryRow(query, postID, commentID).Scan(&postComment.PostID, &postComment.CommentID)
	if err != nil {
		return nil, err
	}
	return &postComment, nil
}

func UpdatePostComment(db *sql.DB, postID, commentID int, updatedPostComment PostComment) error {
	query := `UPDATE post_comments SET post_id=$1, comment_id=$2 WHERE post_id=$3 AND comment_id=$4`
	_, err := db.Exec(query, updatedPostComment.PostID, updatedPostComment.CommentID, postID, commentID)
	return err
}

func DeletePostComment(db *sql.DB, postID, commentID int) error {
	query := `DELETE FROM post_comments WHERE post_id=$1 AND comment_id=$2`
	_, err := db.Exec(query, postID, commentID)
	return err
}

func GetPostComments(db *sql.DB, postID int) ([]Comment, error) {
	query := `
		SELECT c.id, c.content, c.reg_date, u.id AS user_id, u.username, u.avatar
		FROM comments c
		JOIN post_comments pc ON c.id = pc.comment_id
		JOIN users u ON c.id = u.id
		WHERE pc.post_id = $1
		ORDER BY c.reg_date DESC`

	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.RegDate, &comment.UserID, &comment.Username, &comment.Avatar); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
