# Paths to Docker Compose files
LOCAL_COMPOSE_FILE := "./docker-compose.yml"

# Docker Compose command function
define dc_command
docker compose -f $(1) $(2)
endef

# Ensure Docker networks are created
create_local_network:
	docker network inspect microservices_nginx_network >/dev/null 2>&1 || docker network create microservices_nginx_network

# Common commands
up: create_local_network
	$(call dc_command,$(LOCAL_COMPOSE_FILE),up -d --build)

down:
	$(call dc_command,$(LOCAL_COMPOSE_FILE),down -v --remove-orphans)

logs:
	$(call dc_command,$(LOCAL_COMPOSE_FILE),logs -f)

ps:
	$(call dc_command,$(LOCAL_COMPOSE_FILE),ps)


.PHONY: up down logs ps