#!/bin/bash
#set -x
oldfile=${1:?usage: exdif.sh  oldfile newfile [summary]}
newfile=${2:?usage: exdif.sh oldfile newfile [summary]}
summary=${3:-frd}
wc -l $newfile >/dev/null || exit
wc -l $oldfile >/dev/null || exit

IDF1=$summary.tmp.OldId
IDF2=$summary.tmp.NewId
UNIQ=$summary.tmp.UniqueId

cut -f1 $oldfile | sort > $IDF1
cut -f1 $newfile | sort > $IDF2

sort -m $IDF1 $IDF2 |
    uniq -u > $UNIQ |
    sed 's/$/\t/'

grep -F -f$UNIQ $oldfile > $summary-del.txt
grep -F -f$UNIQ $newfile > $summary-add.txt

echo $summary-add ===============
head -10 $summary-add.txt
echo $summary-del ===============
head -10 $summary-del.txt

rm $IDF1
rm $IDF2
rm $UNIQ
