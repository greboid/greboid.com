#!/bin/bash
for file in $(find "/app/static/images/$dir" -name '*.jpg' -o -name '*.png' -o -name '*.jpeg'); do
  cwebp -quiet -m 6 -mt -o "$file.webp" -- "$file"
done