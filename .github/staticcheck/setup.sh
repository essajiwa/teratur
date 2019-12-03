#!/bin/sh

set -e

RELEASES_URL="https://github.com/dominikh/go-tools/releases"

downloadJSON() {
    url="$2"

    echo "Fetching $url.."
    if test -x "$(command -v curl)"; then
        response=$(curl -s -L -w 'HTTPSTATUS:%{http_code}' -H 'Accept: application/json' "$url")
        body=$(echo "$response" | sed -e 's/HTTPSTATUS\:.*//g')
        code=$(echo "$response" | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    elif test -x "$(command -v wget)"; then
        temp=$(mktemp)
        body=$(wget -q --header='Accept: application/json' -O - --server-response "$url" 2> "$temp")
        code=$(awk '/^  HTTP/{print $2}' < "$temp" | tail -1)
        rm "$temp"
    else
        echo "Neither curl nor wget was available to perform http requests."
        exit 1
    fi
    if [ "$code" != 200 ]; then
        echo "Request failed with code $code"
        exit 1
    fi

    eval "$1='$body'"
}

downloadFile() {
    url="$1"
    destination="$2"

    echo "Fetching $url.."
    if test -x "$(command -v curl)"; then
        code=$(curl -s -w '%{http_code}' -L "$url" -o "$destination")
    elif test -x "$(command -v wget)"; then
        code=$(wget -q -O "$destination" --server-response "$url" 2>&1 | awk '/^  HTTP/{print $2}' | tail -1)
    else
        echo "Neither curl nor wget was available to perform http requests."
        exit 1
    fi

    if [ "$code" != 200 ]; then
        echo "Request failed with code $code"
        exit 1
    fi
}

findGoBinDirectory() {
    EFFECTIVE_GOPATH=$(go env GOPATH)
    # CYGWIN: Convert Windows-style path into sh-compatible path
    if [ "$OS_CYGWIN" = "1" ]; then
	EFFECTIVE_GOPATH=$(cygpath "$EFFECTIVE_GOPATH")
    fi
    if [ -z "$EFFECTIVE_GOPATH" ]; then
        echo "Installation could not determine your \$GOPATH."
        exit 1
    fi
    if [ -z "$GOBIN" ]; then
        GOBIN=$(echo "${EFFECTIVE_GOPATH%%:*}/bin" | sed s#//*#/#g)
    fi
    if [ ! -d "$GOBIN" ]; then
        echo "Installation requires your GOBIN directory $GOBIN to exist. Please create it."
        exit 1
    fi
    eval "$1='$GOBIN'"
}

initArch() {
    ARCH=$(uname -m)
    if [ -n "$DEP_ARCH" ]; then
        echo "Using DEP_ARCH"
        ARCH="$DEP_ARCH"
    fi
    case $ARCH in
        amd64) ARCH="amd64";;
        x86_64) ARCH="amd64";;
        i386) ARCH="386";;
        ppc64) ARCH="ppc64";;
        ppc64le) ARCH="ppc64le";;
        s390x) ARCH="s390x";;
        armv6*) ARCH="arm";;
        armv7*) ARCH="arm";;
        aarch64) ARCH="arm64";;
        *) echo "Architecture ${ARCH} is not supported by this installation script"; exit 1;;
    esac
    echo "ARCH = $ARCH"
}

initOS() {
    OS=$(uname | tr '[:upper:]' '[:lower:]')
    OS_CYGWIN=0
    if [ -n "$DEP_OS" ]; then
        echo "Using DEP_OS"
        OS="$DEP_OS"
    fi
    case "$OS" in
        darwin) OS='darwin';;
        linux) OS='linux';;
        freebsd) OS='freebsd';;
        mingw*) OS='windows';;
        msys*) OS='windows';;
	cygwin*)
	    OS='windows'
	    OS_CYGWIN=1
	    ;;
        *) echo "OS ${OS} is not supported by this installation script"; exit 1;;
    esac
    echo "OS = $OS"
}

# identify platform based on uname output
initArch
initOS

# determine install directory if required
if [ -z "$INSTALL_DIRECTORY" ]; then
    findGoBinDirectory INSTALL_DIRECTORY
fi
echo "Will install into $INSTALL_DIRECTORY"

# assemble expected release artifact name
if [ "${OS}" != "linux" ] && { [ "${ARCH}" = "ppc64" ] || [ "${ARCH}" = "ppc64le" ];}; then
    # ppc64 and ppc64le are only supported on Linux.
    echo "${OS}-${ARCH} is not supported by this instalation script"
else
    BINARY="staticcheck_${OS}_${ARCH}.tar.gz"
fi

# if DEP_RELEASE_TAG was not provided, assume latest
if [ -z "$RELEASE_TAG" ]; then
    downloadJSON LATEST_RELEASE "$RELEASES_URL/latest"
    RELEASE_TAG=$(echo "${LATEST_RELEASE}" | tr -s '\n' ' ' | sed 's/.*"tag_name":"//' | sed 's/".*//' )
fi
echo "Release Tag = $RELEASE_TAG"

# fetch the real release data to make sure it exists before we attempt a download
downloadJSON RELEASE_DATA "$RELEASES_URL/tag/$RELEASE_TAG"

BINARY_URL="$RELEASES_URL/download/$RELEASE_TAG/$BINARY"
DOWNLOAD_FILE=$(mktemp)

downloadFile "$BINARY_URL" "$DOWNLOAD_FILE"

INSTALL_NAME="staticcheck"

if [ "$OS" = "windows" ]; then
    INSTALL_NAME="$INSTALL_NAME.exe"
fi

echo "Extraction"
mv "$DOWNLOAD_FILE" "/tmp/$BINARY"
tar -xvzf "/tmp/$BINARY"

echo "Moving executable to $INSTALL_DIRECTORY/$INSTALL_NAME"
mv "staticcheck/$INSTALL_NAME" "$INSTALL_DIRECTORY/$INSTALL_NAME"

echo "Setting executable permissions."
chmod +x "$INSTALL_DIRECTORY/$INSTALL_NAME"
