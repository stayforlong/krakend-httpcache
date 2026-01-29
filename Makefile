mocks: ## Generate mocks for testing
	@echo "generating mock files ..."
	find $(CURDIR) -name "mock_*.go" -not -path "$(CURDIR)/vendor/*" -delete
	go generate -run="mockgen" ./...
	@echo "... done"