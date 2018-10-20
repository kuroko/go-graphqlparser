default: help

.PHONY: help
help: ## Show this help.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Unit test the go source files.
	./scripts/test.sh

.PHONY: generate-lists
generate-lists: ast/lists.go graphql/lists.go ## Generate the linked list types go source files.

.PHONY: generate-walkerevents
generate-walkerevents: ## Generate the walker events go source files.
	./scripts/walker_events.sh

ast/lists.go:
	./scripts/ast-lists.sh

graphql/lists.go:
	./scripts/graphql-lists.sh
