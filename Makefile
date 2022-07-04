install-tools:
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
