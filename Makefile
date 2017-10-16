# 2017-05-28 adbr

all:
	@echo "targets: release, tag, clean"

release:
	sh release.sh

tags:
	uctags -e --extra=q -R .

clean:
	rm -f pogoda
	rm -rf build
	rm -rf release
	rm -f TAGS
