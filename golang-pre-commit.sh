#!/bin/sh

# Script does not handle file names that contain spaces.
# It runs all the checks. If any of the checks fails, it doesn't stop, but rather collects all the errors.

has_errors=0

# get go files to check against
gofiles=$(git diff --cached --name-only --diff-filter=ACM | grep '.go$')
[ -z "$gofiles" ] && exit 0

# go fmt
unformatted=$(gofmt -l $gofiles)
if [ -n "$unformatted" ]; then
	echo >&2 "Go fmt check results:\nFiles must be formatted with gofmt. Please, run:"
	for f in $unformatted; do
		echo >&2 " gofmt -w $PWD/$f"
	done
	echo "\n"
	has_errors=1
fi

# go lint
if golint >/dev/null 2>&1; then # check if golint is installed
	lint_errors=false
	for file in $gofiles ; do
		lint_result="$(golint $file)" # run golint
		if test -n "$lint_result" ; then
			echo "Go lint check result for '$file':\n$lint_result"
			lint_errors=true
			has_errors=1
		fi
	done
	if [ $lint_errors = true ] ; then
		echo "\n"
	fi
else
	echo 'Error: golint is not installed. To install run: "go get -u github.com/golang/lint/golint"' >&2
	exit 1
fi

# go vet
show_vet_header=true
for file in $gofiles ; do
	vet=$(go tool vet $file 2>&1)
	if [ -n "$vet" -a $show_vet_header = true ] ; then
		echo "Go vet check results:"
		show_vet_header=false
	fi
	if [ -n "$vet" ] ; then
		echo "$vet\n"
		has_errors=1
	fi
done

exit $has_errors