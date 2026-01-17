.PHONY: release-patch release-minor release-major tag-patch tag-minor tag-major current-version help

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

tag-patch:
	@echo "Creating tag $(NEXT_PATCH)..."
	git tag -a $(NEXT_PATCH) -m "Release $(NEXT_PATCH)"

tag-minor:
	@echo "Creating tag $(NEXT_MINOR)..."
	git tag -a $(NEXT_MINOR) -m "Release $(NEXT_MINOR)"

tag-major:
	@echo "Creating tag $(NEXT_MAJOR)..."
	git tag -a $(NEXT_MAJOR) -m "Release $(NEXT_MAJOR)"

release-patch: tag-patch
	@echo "Pushing $(NEXT_PATCH) to origin..."
	git push origin $(NEXT_PATCH)
	@echo "Released $(NEXT_PATCH)"

release-minor: tag-minor
	@echo "Pushing $(NEXT_MINOR) to origin..."
	git push origin $(NEXT_MINOR)
	@echo "Released $(NEXT_MINOR)"

release-major: tag-major
	@echo "Pushing $(NEXT_MAJOR) to origin..."
	git push origin $(NEXT_MAJOR)
	@echo "Released $(NEXT_MAJOR)"
