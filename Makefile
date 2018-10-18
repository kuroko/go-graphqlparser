default: help

.PHONY: help

help: ## Show this help.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'

generate-lists: ## Generate the linked list types go source files.
	./scripts/lists.sh
.PHONY: generate-lists

generate-walkerevents: ## Generate the walker events go source files.
	./scripts/walker_events.sh
.PHONY: generate-walkerevents
