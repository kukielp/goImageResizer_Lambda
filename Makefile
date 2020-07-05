.PHONY: build

build:
	sam build
	cp error.jpg .aws-sam/build/goImageResizerFunction/.
	@echo "Build Good, run make deploy"
deploy:
	sam deploy