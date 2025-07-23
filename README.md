# plex-cache

If a tv series episode plays it caches the next 4 episodes on a seperate drive and then using mergerfs it prefers files on cache drive instead of main storage drive.

Basic improvements

- Keep track of diskspace for /cache drive.
- Remove files from /cache.
- Do not cache untill playing last of the cached episode.
- Do not cache if it is the first episode of season one.

## Setup:

**Mount point**

- main drive is mounted as /mnt/drive1
- cache drive is mounted as /mnt/cache
- plex uses /mnt/media

/mnt/media is connected to /mnt/drive1 and using mergerfs if files exists in /mnt/cache it is prefered over /mnt/drive1

**Mergerfs setup**

mergerfs -o defaults,allow_other,use_ino,cache.files=off,category.create=epmfs /mnt/cache=RO:/mnt/drive1 /mnt/media

/mnt/cache=RO:/mnt/drive1 /mnt/media fuse.mergerfs defaults,allow_other,use_ino,cache.files=off,category.create=epmfs 0 0
