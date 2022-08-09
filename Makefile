package = github.com/Crosse/font-install

default: release

define build
	@env GOOS=$(1) GOARCH=$(2) make release/font-install-$(1)-$(2)$(3)
endef

release/font-install-%:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o "$@" $(package)

release/font-install-darwin-universal: release/font-install-darwin-amd64 release/font-install-darwin-arm64
	$(RM) "$@"
	lipo -create -o "$@" $^

.PHONY: release
release:
	mkdir -p release
	$(call build,linux,arm)
	$(call build,linux,amd64)
	$(call build,linux,arm64)

	$(call build,darwin,amd64)
	$(call build,darwin,arm64)
	@make release/font-install-darwin-universal

	$(call build,openbsd,arm)
	$(call build,openbsd,amd64)
	$(call build,openbsd,arm64)

	$(call build,freebsd,arm)
	$(call build,freebsd,amd64)
	$(call build,freebsd,arm64)

	$(call build,windows,amd64,.exe)
	$(call build,windows,arm64,.exe)

.PHONY: zip
zip: release
	find release -type f ! -name '*.zip' -execdir zip -9 "{}.zip" "{}" \;

.PHONY: clean
clean:
	$(RM) -r release
