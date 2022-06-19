#!/usr/bin/env bash
set -e

HABR_USER_ID="$1"

if [[ -z "$HABR_USER_ID" ]]; then
	echo "Usage: $0 <user_id>" >&2
	exit 1
fi

echo "Building docker image" >&2
DOCKER_IMAGE=$(docker build . -q)

echo "" >&2
echo "Fetching new articles" >&2
mkdir -p "./out/$HABR_USER_ID/html"
docker run -t --rm -v "$(pwd)/out/$HABR_USER_ID/html:/mnt" $DOCKER_IMAGE --user-id "$HABR_USER_ID" --output /mnt

echo "" >&2
echo "Converting articles to EPUB" >&2
find "./out/$HABR_USER_ID/html" -maxdepth 1 -type d -name "*" | while read dir; do
	dirname=$(basename "$dir")

	find "$dir" -maxdepth 1 -type f -name "*.html" | while read file; do
		filename=$(basename "$file" ".html")

		if [[ -f "./out/$HABR_USER_ID/epub/$dirname/$filename.epub" ]]; then
			continue
		fi

		echo "- processing \"$filename\"" >&2

		cp "$file" "./out/$HABR_USER_ID/html/$dirname/in.html"

		OUTPUT=$(docker run -t --rm -v "$(pwd)/out/$HABR_USER_ID/html/$dirname:/mnt" larrycai/ebook-convert:latest ebook-convert \
			"/mnt/in.html" "/mnt/out.epub" --no-svg-cover --embed-all-fonts)

		mkdir -p "./out/$HABR_USER_ID/epub/$dirname"
		rm "./out/$HABR_USER_ID/html/$dirname/in.html"
		mv "./out/$HABR_USER_ID/html/$dirname/out.epub" "./out/$HABR_USER_ID/epub/$dirname/$filename.epub"
	done
done
