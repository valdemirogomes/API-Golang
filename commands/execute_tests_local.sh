#!/bin/zsh

COVERAGE_REPORT=/tmp/coverage.out
TMP_COVIGNORE=/tmp/.covignore

function parse_exclusions {
  echo -e "\nRunning exclusions for file $1..."
  echo -e "$(cat $1)\n" > $TMP_COVIGNORE

  while read LINE; do
    for FILE in $LINE; do
      if [ ! -z "$FILE" ]; then
        echo "- Excluding $LINE..."

        RESULT=$(cat $COVERAGE_REPORT | grep -v $FILE)
        echo $RESULT > $COVERAGE_REPORT
      fi
    done
  done < $TMP_COVIGNORE
}

function show_help {
  echo -e 'Go Test Utils\n'
  echo -e '\tgotest [-v] [-open]\n'
  echo '  -v\tEnables verbose output for tests'
  echo '  -open\tCalls go tool cover to show HTML report'
  exit 0
}

ARGS=""
OPEN=false
while [ "$1" != "" ]; do
  case $1 in
    open | -open ) OPEN=true
    ;;
    -v | -verbose | --verbose | verbose | v ) ARGS="-v"
    ;;
    -h | --help ) show_help
  esac
  shift
done

go test ./... -covermode=atomic -coverprofile=$COVERAGE_REPORT -coverpkg=./... -count=1 $ARGS

if [ -f .covignore ]; then
  parse_exclusions .covignore
fi

if [ -f src/api/.covignore ]; then
  parse_exclusions src/api/.covignore
fi

COVERAGE=$(go tool cover -func=/tmp/coverage.out | grep total: | awk '{print $3}')
echo -e "\nTOTAL COVERAGE: $COVERAGE"

if [ $OPEN = true ]; then
  go tool cover -html=$COVERAGE_REPORT
fi