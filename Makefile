CURRENT_TIME := `date +%s`

# make create-migration NAME="create_users_table"
create-migration:
	@[ ! -z ${NAME} ]
	mkdir -p assets/migrate
	python scripts/makefile_helper/helper.py write_migration ${NAME}
	@go fmt assets/migrate/*

test:
	@echo "=================================================================================="
	@echo "Coverage Test"
	@echo "=================================================================================="
	go fmt ./... && go test -coverprofile coverage.cov -cover ./... # use -v for verbose
	@echo "\n"
	@echo "=================================================================================="
	@echo "All Package Coverage"
	@echo "=================================================================================="
	go tool cover -func coverage.cov