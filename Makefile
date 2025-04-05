swag-gen:
	swag init -g internal/http/comment.go -o docs --parseDependency --parseInternal
	swag init -g internal/http/post.go -o docs --parseDependency --parseInternal
	swag init -g internal/http/user.go -o docs --parseDependency --parseInternal
