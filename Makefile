.PHONY: serve
serve: content
	go run . -cert=./testdata/crt.pem -key=./testdata/key.pem


.PHONY: content
content: HelloWorld.jar
	rm -rf $@ || true
	mvn install:install-file \
		-DgroupId=carto.run \
		-DartifactId=hello-world \
		-Dversion=0.0.1 \
		-Dfile=$^ \
		-Dpackaging=jar \
		-DgeneratePom=true \
		-DlocalRepositoryPath=$@ \
		-DcreateChecksum=true

%.jar: %.java
	javac $^
	jar cvfe $@ $* ./*.class


.PHONY: release
release: content
	mkdir -p release
	ko resolve --bare -f ./config > ./release/release.yaml
