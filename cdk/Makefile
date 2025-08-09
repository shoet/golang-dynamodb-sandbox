.DEFAULT_GOAL := help

deploy: ## deploy by cdk
	npx cdk deploy --require-approval never GolangDynamoDBSandbox

deploy-local: ## deploy by cdklocal
	npx cdklocal destroy --force GolangDynamoDBSandbox-local
	npx cdklocal deploy GolangDynamoDBSandbox-local

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
    awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
