# 2017-05-28 adbr

all:
	@echo "targets: release"

release:
	sh release.sh

clean:
	rm -rf build
	rm -rf release
