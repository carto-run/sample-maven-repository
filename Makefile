.PHONY: serve
serve: content
	go run main.go -addr=:8080


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
	jar cvf $@ $*.class


.PHONY: release
release: content
	mkdir -p release
	ko resolve --bare -f ./config > ./release/release.yaml
