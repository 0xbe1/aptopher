.PHONY: release-patch release-minor release-major current-version help

# Get the latest version tag, default to v0.0.0 if none exists
CURRENT_VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
VERSION_PARTS := $(subst ., ,$(subst v,,$(CURRENT_VERSION)))
MAJOR := $(word 1,$(VERSION_PARTS))
MINOR := $(word 2,$(VERSION_PARTS))
PATCH := $(word 3,$(VERSION_PARTS))

# Calculate next versions
NEXT_PATCH := v$(MAJOR).$(MINOR).$(shell echo $$(($(PATCH)+1)))
NEXT_MINOR := v$(MAJOR).$(shell echo $$(($(MINOR)+1))).0
NEXT_MAJOR := v$(shell echo $$(($(MAJOR)+1))).0.0

help:
	@echo "Release targets:"
	@echo "  make release-patch  - Release a patch version (bug fixes)"
	@echo "  make release-minor  - Release a minor version (new features)"
	@echo "  make release-major  - Release a major version (breaking changes)"
	@echo ""
	@echo "Current version: $(CURRENT_VERSION)"
	@echo "  Next patch: $(NEXT_PATCH)"
	@echo "  Next minor: $(NEXT_MINOR)"
	@echo "  Next major: $(NEXT_MAJOR)"

current-version:
	@echo $(CURRENT_VERSION)

release-patch:
	@echo "Releasing $(NEXT_PATCH)..."
	git tag $(NEXT_PATCH)
	git push origin $(NEXT_PATCH)
	gh release create $(NEXT_PATCH) --generate-notes --edit
	@echo "Released $(NEXT_PATCH)"

release-minor:
	@echo "Releasing $(NEXT_MINOR)..."
	git tag $(NEXT_MINOR)
	git push origin $(NEXT_MINOR)
	gh release create $(NEXT_MINOR) --generate-notes --edit
	@echo "Released $(NEXT_MINOR)"

release-major:
	@echo "Releasing $(NEXT_MAJOR)..."
	git tag $(NEXT_MAJOR)
	git push origin $(NEXT_MAJOR)
	gh release create $(NEXT_MAJOR) --generate-notes --edit
	@echo "Released $(NEXT_MAJOR)"
