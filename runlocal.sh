#! /usr/bin/bash
set -e

rootdir=/tmp/video-in-be
if [ -d $rootdir ]; then
    sudo rm -rf $rootdir
fi
mkdir -p $rootdir
export VIDEOIN_UNCLAIMEDDIR=$rootdir/unclaimed
mkdir -p $VIDEOIN_UNCLAIMEDDIR
mkdir -p $VIDEOIN_UNCLAIMEDDIR/disc1 $VIDEOIN_UNCLAIMEDDIR/disc2 $VIDEOIN_UNCLAIMEDDIR/disc3 $VIDEOIN_UNCLAIMEDDIR/disc4
cp testdata/testdata_sample_640x360.mkv $VIDEOIN_UNCLAIMEDDIR/disc1
cp testdata/testdata_sample_640x360.mkv $VIDEOIN_UNCLAIMEDDIR/disc2
cp testdata/testdata_sample_640x360.mkv $VIDEOIN_UNCLAIMEDDIR/disc3
cp testdata/testdata_sample_640x360.mkv $VIDEOIN_UNCLAIMEDDIR/disc4
export VIDEOIN_STATEDIR=$rootdir/state
mkdir -p $VIDEOIN_STATEDIR
export VIDEOIN_PROJECTDIR=$rootdir/project
mkdir -p $VIDEOIN_PROJECTDIR
export VIDEOIN_THUMBSDIR=$rootdir/thumbs
mkdir -p $VIDEOIN_THUMBSDIR
export VIDEOIN_LIBRARYDIR=$rootdir/library
mkdir -p $VIDEOIN_LIBRARYDIR
source tmdb_key.sh

docker build -t video-in-be .
docker run --rm \
    -v ${rootdir}:${rootdir} \
    -v /nas/media:/nas/media \
    --env VIDEOIN_TMDBKEY \
    --env VIDEOIN_UNCLAIMEDDIR \
    --env VIDEOIN_STATEDIR \
    --env VIDEOIN_PROJECTDIR \
    --env VIDEOIN_THUMBSDIR \
    --env VIDEOIN_LIBRARYDIR \
    -it -p 25004:25004 --name video-in-be video-in-be "${@}"