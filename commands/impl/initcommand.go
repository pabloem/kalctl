package impl

const AUTOCOMPLETE_SCRIPT = `#!/bin/bash
_kalctl()
{
  # echo "COMP_LINE: <$COMP_LINE>"
  reply_words=$($COMP_LINE --help --short)
  cur=""

  # If COMP_LINE does not contain a full command, we autocomplete
  # that command (remove the last word and run again).
  if [ $? -ne 0 ] || [[ " ${reply_words[@]} " =~ "Error" ]]; then
    COMP_LINE=$(echo $COMP_LINE | rev | cut -d " " -f 2- | rev)
    reply_words=$($COMP_LINE --help --short)
    cur="${COMP_WORDS[COMP_CWORD]}"
  fi
  COMPREPLY=( $(compgen -W "${reply_words}" -- ${cur}) )

  return 0
}
complete -F _kalctl kalctl
`
