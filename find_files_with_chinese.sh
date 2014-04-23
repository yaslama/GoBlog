
# Files with non-ascii characters:
# (install package ack or ack-grep or similar..)
ack -l "[\x80-\xFF]"

echo -e "\n\n"

# the matching lines:
ack --no-color "[\x80-\xFF]"

