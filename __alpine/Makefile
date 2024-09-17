all: 
	cd ../src
	CGO_ENABLED=0 go build . -o ./prometheusSDsendHost.exe

prepare:
	sudo /opt/bin/install_golang.sh 1.23.1 amd64
	sudo apk update --no-cache

install:
	test -d $(DESTDIR)$(PREFIX)/bin || mkdir -p $(DESTDIR)$(PREFIX)/bin
	install -pm 755 $(BUILDDIR)/prometheusSDsendHost.exe $(DESTDIR)$(PREFIX)/bin/prometheusSDsendHost
	strip $(DESTDIR)$(PREFIX)/bin/prometheusSDsendHost
