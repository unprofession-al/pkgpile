VERSION=`git tag -l --points-at HEAD`
COMMITID=`git rev-parse HEAD`
go install -ldflags "-X main.version=${VERSION} -X main.commitId=${COMMITID}"
