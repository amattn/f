#!/bin/sh

# Last edit 2020-08-14

set -o nounset
set -o errexit

# you can add the following lines to .git/hooks/pre-commit to auto bump build num:
#
# ./version_bump.sh
# git add version.go

VERSION_GO_FILENAME="version.go"

usage(){
  echo "Usage: $0 precommit|major|minor|patch|build|tag"
  echo "  $0 build (just update buildnum, default)"
  echo "  $0 patch (1.0.PATCH)"
  echo "  $0 minor (1.MINOR.0)"
  echo "  $0 major (MAJOR.0.0)"
  echo "  $0 tag (git tag with the X.Y.Z version number)"
  echo "  $0 precommit (append the following lines to .git/hooks/pre-commit)"
  echo "               ./version_bump.sh"
  echo "               git add version.go"
}

bump_build(){
  if ! CURRENT_BUILD_NUM=$(grep -o "internalBuildNumber\\s*\(int64\|int\)\?\\s*=\\s*[0-9]\\+" $VERSION_GO_FILENAME | grep -o " = [0-9]\\+" | grep -o "[0-9]\\+") ; then
    echo "Fatal Error: Cannot find/parse internalBuildNumber"
    exit 16
  fi

  NEW_NUM=$((CURRENT_BUILD_NUM+1))
  echo "Bumping build number, from ${CURRENT_BUILD_NUM} to $NEW_NUM"

  sed -i.bak "s/internalBuildNumber[[:space:]]*[int64]*[[:space:]]*=[[:space:]]*[0-9]*/internalBuildNumber\\ int64\\ =\\ ${NEW_NUM}/g" $VERSION_GO_FILENAME


  NEW_TS=$(date +%s)
  echo "Bumping build timestamp new=$NEW_TS"
  sed -i.bak "s/internalBuildTimestamp[[:space:]]*[int64]*[[:space:]]*=[[:space:]]*[0-9][0-9]*/internalBuildTimestamp\\ int64\\ =\\ ${NEW_TS}/g" $VERSION_GO_FILENAME
  rm -f version.go.bak

  # cleanup
  go fmt version.go > /dev/null

  FINAL_VERSION_STRING=$(grep -o "internalVersionString\\s*=\\s\"[vV]*[0-9]*\.[0-9]*\.[0-9]*\"\\+" $VERSION_GO_FILENAME | grep -o "[0-9]*\.[0-9]*\.[0-9]*")
  echo "v${FINAL_VERSION_STRING}"
}

bump_version_string(){
  echo "Bumping ${POSITION} number, from ${XYZ_VERSION_STRING} to ${NEW_VERSION_STRING}"
  sed -i.bak "s/internalVersionString[[:space:]]*=[[:space:]]*\"[vV]*[0-9]*\.[0-9]*\.[0-9]*\"/internalVersionString\\ =\\ \"${NEW_VERSION_STRING}\"/g" $VERSION_GO_FILENAME

  bump_build
}


bump_patch(){
  echo "patching"
  NEW_PATCH_NUM=$((PATCH_NUM+1))
  NEW_VERSION_STRING=${MAJOR_NUM}.${MINOR_NUM}.${NEW_PATCH_NUM}
  bump_version_string
}

bump_minor(){
  NEW_MINOR_NUM=$((MINOR_NUM+1))
  NEW_PATCH_NUM="0"
  NEW_VERSION_STRING=${MAJOR_NUM}.${NEW_MINOR_NUM}.${NEW_PATCH_NUM}
  bump_version_string
}

bump_major(){
  NEW_MAJOR_NUM=$((MAJOR_NUM+1))
  NEW_MINOR_NUM="0"
  NEW_PATCH_NUM="0"
  NEW_VERSION_STRING=${NEW_MAJOR_NUM}.${NEW_MINOR_NUM}.${NEW_PATCH_NUM}
  bump_version_string
}

set_tag(){
  FINAL_VERSION_STRING=$(grep -o "internalVersionString\\s*=\\s\"[vV]*[0-9]*\.[0-9]*\.[0-9]*\"\\+" $VERSION_GO_FILENAME | grep -o "[0-9]*\.[0-9]*\.[0-9]*")
  echo "set_tag: v$FINAL_VERSION_STRING"
  git add version.go
  echo "git commit -m \"tagging version v$FINAL_VERSION_STRING\""
  if ! git commit -m "\"tagging version v$FINAL_VERSION_STRING\"" ; then
    echo "git commit failure.  Did you forget to bump major, minor or patch?"
    exit 18
  fi

  echo "git tag \"v$FINAL_VERSION_STRING\""
  git tag "v$FINAL_VERSION_STRING"

  echo ""
  echo "if you need to remove the tag locally, use:"
  echo "git tag -d v$FINAL_VERSION_STRING"
  echo ""
  echo "use the following command to push tag to origin:"
  echo "git push && git push --tags"
}

precommit(){
  #TODO need to check if these lines already exist in the file
  {
    echo "#!/bin/sh"
    echo ""
    echo '# automatically bump build number on commit'
    echo './version_bump.sh'
    echo 'git add version.go'
  } >> .git/hooks/pre-commit

  chmod +x .git/hooks/pre-commit
}


if [ $# -eq 0 ]
then
  bump_build
  exit 0
elif [ $# -gt 1 ]
then
  echo "too many arguments supplied"
  usage
  exit 1
fi


while [ "$1" != "" ]; do

  PARAM=$(echo "$1" | awk -F= '{print $1}')

  # should be of form x.y.z
  XYZ_VERSION_STRING=$(grep -o "internalVersionString\\s*=\\s\"[vV]*[0-9]*\.[0-9]*\.[0-9]*\"\\+" "$VERSION_GO_FILENAME" | grep -o "[0-9]*\.[0-9]*\.[0-9]*")

  MAJOR_NUM=$(echo "$XYZ_VERSION_STRING" | cut -d '.' -f 1)
  MINOR_NUM=$(echo "$XYZ_VERSION_STRING" | cut -d '.' -f 2)
  PATCH_NUM=$(echo "$XYZ_VERSION_STRING" | cut -d '.' -f 3)

  POSITION=$1

  case $PARAM in
      -h | --help)
          usage
          exit 0
          ;;
      build)
        bump_build
        exit 0
        ;;
      patch)
        bump_patch
        exit 0
        ;;
      minor)
        bump_minor
        exit 0
        ;;
      major)
        bump_major
        exit 0
        ;;
      tag)
        set_tag
        exit 0
        ;;
      precommit)
        precommit
        exit 0
        ;;
      *)
          echo "ERROR: unknown parameter \"$PARAM\""
          usage
          exit 1
          ;;
  esac
  shift
done

