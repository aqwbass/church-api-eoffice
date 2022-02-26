
test:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out -o coverage.html
	export unit_total=$$(go test ./... -v  | grep -c RUN) && echo "Unit Test Total: $$unit_total" && export coverage_total=$$(go tool cover -func cover.out | grep total | awk '{print $$3}') && echo "Coverage Total: $$coverage_total"