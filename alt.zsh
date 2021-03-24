#!/usr/bin/env zsh


find . \( -name 'node_modules' -or -name 'es' -or -name 'dist' -or -name 'target' \) -prune \
	-and -type f -name '*.css' -or -name '*.scss' -or -name '*.sass' \
	| xargs grep -E "$1" \
	| sed -E $'s/^([^:]*):/\\1\t/g' \
	| sed -E $'s/([ {]*$)//g' \
	| sed -E $'s/(\t +)/\t/g'

