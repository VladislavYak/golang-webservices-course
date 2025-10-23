package comment

type CommentRepository interface {
	AddComment(string, *Comment) error
	DeleteComment(string, string) error
}
