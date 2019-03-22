.PHONY: default clean hugo hugo-build

default: hugo

clean:
	rm -rf public/

hugo-build: clean hugo-themes
	hugo --enableGitInfo --source .

hugo:
	hugo server --disableFastRender --enableGitInfo --watch --source .
	# hugo server -D

hugo-themes: hugo-theme-clean
	mkdir -p themes
	git clone --depth=1 https://github.com/bep/docuapi.git themes/docuapi
	rm -rf doc/themes/docuapi/.git themes/docuapi/*.go

hugo-theme-clean:
	rm -rf themes