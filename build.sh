#!/usr/bin/env bash
set -ex
IFS=$' \n\t'

export PROVIDER_VERSION="0.5.0"
export DISTDIR="$PWD/dist"
export WORKDIR="$PWD"

export GOX_MAIN_TEMPLATE="$DISTDIR/{{.OS}}/{{.Dir}}_v${PROVIDER_VERSION}"
export GOX_ARCH="amd64"
export GOX_OS=${*:-"linux darwin"}

# We'll use gox to cross-compile
GO111MODULE=off go get github.com/mitchellh/gox
# We just assume the cross toolchains are already installed, since on Debian
# there are deb packages for those.

# Build the provider
gox -arch="$GOX_ARCH" \
    -os="$GOX_OS" \
    -ldflags="-X github.com/saymedia/terraform-buildkite/buildkite/version.Version=${PROVIDER_VERSION}" \
    -output="$GOX_MAIN_TEMPLATE" \
    github.com/saymedia/terraform-buildkite/cmd/terraform-provider-buildkite

# ZZZZZZZZZZZZZZZZZZZZIPPIT
echo "--- Build done"
for os in $GOX_OS; do
    for arch in $GOX_ARCH; do
        echo "--- Zipping $os/$arch"
        cd "$DISTDIR/$os"
        zip ../terraform-provider-buildkite-v"${PROVIDER_VERSION}"-"$os"-"$arch".zip ./*
    done
done
echo "--- DING! Fries are done"
cd "$DISTDIR"
openssl dgst -r -sha256 ./*.zip > sha256s.txt
exit 0
