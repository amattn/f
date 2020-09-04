#!/bin/sh

# if one of our commands returns an error, stop execution of this script
set -o errexit

COMPONENT="f"
GO_COMMAND="go"

## build on the native or default platform
echo "************************"
echo "building native platform"
$GO_COMMAND build

# update build timestamp
VERSION_GO_FILENAME="version.go"
NEW_TS=$(date +%s)
echo "Bumping build timestamp"
sed -i.bak "s/internalBuildTimestamp[[:space:]]*[int64]*[[:space:]]*=[[:space:]]*[0-9][0-9]*/internalBuildTimestamp\\ int64\\ =\\ ${NEW_TS}/g" $VERSION_GO_FILENAME
rm -f version.go.bak
# cleanup
go fmt version.go > /dev/null


echo "************************"
echo "crosscompiling: gox -osarch=\"darwin/amd64\" -osarch=\"linux/arm\""
gox -osarch="darwin/amd64" -osarch="linux/arm"

# should be of form x.y.z
XYZ_VERSION_STRING=$(grep -o "internalVersionString\s*=\s\"[vV]*.*\"" $VERSION_GO_FILENAME | perl -nle'print $& while m{(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?}g')
SUFFIX="v$XYZ_VERSION_STRING"

mv "${COMPONENT}_darwin_amd64" "${COMPONENT}_darwin_amd64_$SUFFIX"
mv "${COMPONENT}_linux_arm" "${COMPONENT}_linux_arm_$SUFFIX"

# display build metadata
if ! CURRENT_BUILD_NUM=$(grep -o "internalBuildNumber\\s*\(int64\|int\)\?\\s*=\\s*[0-9]\\+" $VERSION_GO_FILENAME | grep -o " = [0-9]\\+" | grep -o "[0-9]\\+") ; then
  echo "Error: Cannot find/parse internalBuildNumber"
fi
echo ""
echo "internalBuildTimestamp $NEW_TS"
echo "   internalBuildNumber $CURRENT_BUILD_NUM"
echo " internalVersionString v${XYZ_VERSION_STRING}"
