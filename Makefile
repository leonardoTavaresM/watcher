IMAGE_NAME=watcher
TAG=dev

# Caminho local dos sources
DEV_PATH=/home/leonardomalt/Documents/dev

# Build da imagem
build:
	docker build -t $(IMAGE_NAME):$(TAG) .

# Build sem cache
build-nc:
	docker build --no-cache -t $(IMAGE_NAME):$(TAG) .


# Rodar o container em modo interativo com watcher ativo
run:
	docker run --rm -it \
		-e WATCH_PATH=/app/dev \
		-v $(DEV_PATH):/app/dev \
		$(IMAGE_NAME):$(TAG)

# Entrar no container com shell
shell:
	docker run --rm -it \
		-v $(DEV_PATH):/app/dev \
		$(IMAGE_NAME):$(TAG) sh

# Remover a imagem
clean:
	docker rmi $(IMAGE_NAME):$(TAG) || true

# Ciclo completo: limpar, rebuildar e rodar
rebuild: clean build run