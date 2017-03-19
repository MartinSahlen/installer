add_delimiter(){
  l=$(cat $1)
  printf "%s\nSPLIT_HERE_PLEASE\n" "${l}"
  return 0
}
export -f add_delimiter

CASK_PATH=/usr/local/Homebrew/Library/Taps/caskroom/homebrew-cask/
find $CASK_PATH -type d -name "Casks" -exec ls -1 {} \; | xargs printf "${CASK_PATH}Casks/%s \n" | xargs -n 1 -I CASK bash -c 'add_delimiter CASK' > brew-cask.txt

brew info --json=v1 --all > brew.json
