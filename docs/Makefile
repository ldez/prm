.PHONY: default clean hugo hugo-build

default: hugo

clean:
	rm -rf public/ resources/

hugo-build: clean hugo-themes
	hugo --enableGitInfo --source .

hugo:
	hugo server --disableFastRender --enableGitInfo --watch --source .
	# hugo server -D

hugo-themes: hugo-theme-clean
	mkdir -p themes
	git clone https://github.com/bep/docuapi.git themes/docuapi
	cd themes/docuapi; git checkout 673861f6c264ce5639fa2cb8b1ff2306670db549; cd -
	rm -rf themes/docuapi/.git themes/docuapi/*.go

hugo-theme-clean:
	rm -rf themes