.PHONY: all
all: build

.PHONY: build
build: build-binary build-docker-image

.PHONY: build-binary
build-binary: check-compliance
	$(call install_app,"sensormockery","sensormockery/sensormockery.go")

.PHONY: clean
clean: clean-binary clean-docker-image

.PHONY: clean-binary
clean-binary: check-compliance
	$(call uninstall_app,"sensormockery")

.PHONY: check-compliance
check-compliance:
	$(call echo_purple,"Checking environments variables...")
	$(call check_env,GOBIN,$(GOBIN))
	$(call check_env,DOCKER_TAG,$(DOCKER_TAG))
	$(call check_env,CONTAINER_PORT,$(CONTAINER_PORT))
	$(call check_env,HOST_PORT,$(HOST_PORT))

	$(call echo_purple,"Checking commands...")
	$(call check_cmds,go,docker)

# Docker targets
.PHONY: build-docker-image
# Optimized to build smallest possible docker image
build-docker-image: check-compliance
	$(call install_app,"sensormockery","sensormockery/sensormockery.go",true)

	$(call echo_purple,"Building docker image $(DOCKER_TAG)...")
	@docker build --build-arg CONTAINER_PORT -t $(DOCKER_TAG) .
	$(call echo_green,"Successfully built image $(DOCKER_TAG)")

	@rm "sensormockery"
	$(call echo_green,"Successfully removed static sensormockery binary")

.PHONY: clean-docker-image
clean-docker-image: check-compliance
	$(call clean_docker_image)

.PHONY: run-docker-image
run-docker-image:
	$(call echo_purple,"Starting docker image $(DOCKER_TAG)...")
	@docker run -d -p $(HOST_PORT):$(CONTAINER_PORT) $(DOCKER_TAG)
	$(call echo_green,"Successfully started $(DOCKER_TAG)")

.PHONY: kill-docker-image
kill-docker-image:
	$(call echo_purple,"Killing container with tag $(DOCKER_TAG)...")
	@docker ps -q --filter ancestor="$(DOCKER_TAG)" | xargs -I {} docker kill {}
	$(call echo_green,"Killed container with tag $(DOCKER_TAG)...")

.PHONY: push-docker-image
push-docker-image:
ifeq ($(DOCKERHUB_USERNAME),)
	$(call echo_red,"No dockerhub username provided.")
	@exit 1
else ifeq ($(DOCKERHUB_PASSWORD),)
	$(call echo_red,"No dockerhub password provided.")
	@exit 1
endif

	$(call echo_purple,"Logging into DockerHub...")
	@docker login --username=$(DOCKERHUB_USERNAME) --password-stdin <<< $(DOCKERHUB_PASSWORD)

	$(call echo_purple,"Pushing image $(DOCKER_TAG)...")
	@docker push $(DOCKER_TAG)
	$(call echo_green,"Successfully pushed $(DOCKER_TAG)")

.PHONY: run-system-tests
run-system-tests:
	$(call echo_purple,"Running system tests...")
ifeq ($(skip_update), true)
	$(call echo_purple,"Skipping docker image update...")
else
	@make build-docker-image
	@make push-docker-image
endif
	@make run-docker-image
	@ginkgo -r
	@make kill-docker-image

.PHONY: clean-code
clean-code:
	$(call echo_purple,"Cleaning code...")
	go fmt ./...
	go vet ./...
	golint ./...


# Define echo colors
# Use colors as follows:
# - red for error
# - green for success
# - purple for info in progress
define echo_red
	@echo "\033[31m$(1)\033[0m"
endef

define echo_green
	@echo "\033[32m$(1)\033[0m"
endef

define echo_purple
	@echo "\033[95m$(1)\033[0m"
endef

# Check if environment is set-up properly
# 1 - env_var name, 2 - env_var value 
define check_env
	$(if $(2),,
	@# else
		$(call echo_red,"$(1) is not set")
		@exit 1)
endef

define check_cmds
	$(foreach cmd, $(1),
		$(if $(shell command -v $(cmd) &> /dev/null || echo "not found"),
		@# then
			$(call echo_red,"$(cmd) not found")
			@exit 1)
	)
endef

# Install app
# 1 - app_name, 2 - app_path, 3 - here
define install_app
	$(call echo_purple,"Getting go dependencies...")
	@go get -d -v ./...
	$(if $(3),
	@# then
		$(call echo_purple,"Building static sensormockery binary...")
		@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(1) "cmd/$(2)",
	@# else
		$(call echo_purple,"Installing $(1)...")
		@go install "cmd/$(2)"
	)

	$(call echo_green,"Successfully built $(1)")
endef

# Uninstall app
# 1 - app_name
define uninstall_app
	$(if $(shell test -f "$(GOBIN)/$(1)" &> /dev/null && echo "found"),
	@# then
		$(call echo_purple,"Uninstalling $(GOBIN)/$(1)...")
		@rm "$(GOBIN)/$(1)",
	@# else
		$(call echo_purple,"$(GOBIN)/$(1) not found. Skipping..."))
endef

# clean_docker_image is put in a function maintain indempotency
define clean_docker_image
	$(if $(shell docker images -q $(DOCKER_TAG)),
	@# then
		$(call echo_purple,"Removing docker image $(DOCKER_TAG)...")
		@docker image rm $(DOCKER_TAG)
		$(call echo_green,"Successfully removed image $(DOCKER_TAG)"),
	@# else
		$(call echo_purple,"Docker image $(DOCKER_TAG) not found. Skipping..."))
endef