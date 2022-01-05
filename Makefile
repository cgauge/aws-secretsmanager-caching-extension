SHELL=/bin/bash
IMAGE_NAME=customergauge/aws-secretsmanager-caching-sidecar

build-docker-image:
	docker build . -t ${IMAGE_NAME}:extension --target build-extension
	docker build . -t ${IMAGE_NAME}

publish-docker-image: build-docker-image
	(test $(TAG_VERSION)) && echo "Tagging images with \"${TAG_VERSION}\"" || echo "You have to define environment variable TAG_VERSION"
	test $(TAG_VERSION)

	docker image tag ${IMAGE_NAME}:latest ${IMAGE_NAME}:${TAG_VERSION}
	docker image push ${IMAGE_NAME}:latest
	docker image push ${IMAGE_NAME}:${TAG_VERSION}

build-lambda-extension: build-docker-image
	mkdir -p .build/extensions
	rm -f .build/extension.zip
	CID=$$(docker create --entrypoint=scratch ${IMAGE_NAME}:extension) ; \
	docker cp $${CID}:/cache-server .build/extensions
	cd .build; zip -qq -y -r extension.zip extensions

publish-lambda-extension: build-lambda-extension
	declare -a REGIONS=( "eu-west-1" "us-east-1" "ap-southeast-2" ) ; \
	for region in "$${REGIONS[@]}" ; do \
		aws lambda publish-layer-version \
			--region $$region \
			--layer-name 'secretsmanager-caching-extension' \
			--description 'Cache secret server extension' \
			--license-info MIT \
			--output text \
			--query Version \
			--zip-file fileb://.build/extension.zip; \
	done

.PHONY: publish-docker-image build-docker-image publish-lambda-extension build-lambda-extension
