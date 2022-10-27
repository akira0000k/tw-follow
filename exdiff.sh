#!/bin/bash
#
#  exdiff.sh newfile [oldfile]
#  launch  XCode Filemerge.app
#    
#  DIFF=diff exdiff.sh newfile [oldfile]
#  run GNU diff
#  
#  DIFF=no  exdiff.sh newfile [oldfile]
#  create added/deleted file only
#
#set -x
file=${1:?usage: exdif.sh  newfile oldfile}
befile=${2}

wc -l -- "$file" || exit

if [[ "$file" =~ follower ]]; then
    echo compare follower
    type=follower
    summary=flw
else
    echo compare friend
    type=friend
    summary=frd
fi
parentdir=$(dirname "$file")/..

if [ -z "$befile" ]; then
    # when no before-file, guess it
    befile=$(ls  $parentdir/20*/$type*.log | tail -2 | head -1)
fi

wc -l -- "$befile" || exit


# find added/deleted usernames
shelldir=$(dirname $0)
$shelldir/exonly.sh "$befile" "$file" $summary


# show diff
if [ "$DIFF" = no ]; then
    :
elif [ "$DIFF" = diff ]; then
    diff "$befile" "$file"
elif [ "$DIFF" = open ]; then
    # XCode Filemerge commandline interface
    opendiff "$befile" "$file" &
#elif [ "$DIFF" = meld ]; then
else
    meld "$befile" "$file" &
fi
