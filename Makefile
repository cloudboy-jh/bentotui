.PHONY: build test lint release-check release

build:
	go build ./...

test:
	go test ./...

# release-check: build and package locally without uploading anything.
# Use this to verify the release artifacts before pushing a tag.
release-check:
	GITHUB_TOKEN=$(shell gh auth token) goreleaser release --snapshot --clean --skip=publish

# release: DO NOT run this locally for real version tags.
# Push the tag and let CI (release.yml) handle the actual GitHub release.
# Running this locally on a real tag will cause CI to fail with "already_exists".
#
# To release:
#   git tag -a vX.Y.Z -m "vX.Y.Z — <summary>"
#   git push origin vX.Y.Z
#   (CI handles the rest)
release:
	@echo "Do not run 'make release' locally for real version tags."
	@echo "Push the tag and let CI handle the release."
	@echo ""
	@echo "To do a local snapshot check: make release-check"
	@exit 1
