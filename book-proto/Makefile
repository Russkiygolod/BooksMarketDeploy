REPO_NAME=book-proto
IMAGE_VERSION = 1.47_0

.PHONY: generate_go
generate_go:
	@rm -rf go
	@mkdir go
	@docker run -v .:/defs namely/protoc-all:${IMAGE_VERSION} -i schemas -d schemas -o . -l go
	@mv /Users/user/Documents/go Linux/module5/book-proto/Users/user/Documents/go Linux/module5/book-proto/go/api/books/${REPO_NAME}/go .
	@rm -rf /Users/user/Documents/go Linux/module5/book-proto/Users/user/Documents/go Linux/module5/book-proto/go/api/books