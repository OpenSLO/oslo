# load_lib
# ========
#
# Summary: Load a given bats library.
#
# Usage: load_lib <name>
#
# Options:
#   <name>    Name of the library to load.
load_lib() {
  local name="$1"
  load "/usr/lib/bats/${name}/load.bash"
}
