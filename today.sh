#!/bin/bash
set -x

today="${1:-$(date +%Y%m%d)}"

mkdir $today || exit
 
./tw-follow -followers > $today/follower-$today.log
./tw-follow -friends   > $today/friends-$today.log
echo >&2

cd $today
../exdiff.sh follower-*.log

../exdiff.sh friends-*.log

