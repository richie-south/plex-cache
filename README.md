# plex-cache

If a tv series episode starts playing it caches the next 4 episodes on a seperate drive as a cache.

Basic improvements

- Keep track of diskspace for /cache drive.

## Setup:

### Drive mount points

- main drive is mounted as /mnt/drive1
- cache drive is mounted as /mnt/cache
- plex uses /mnt/media

`/mnt/media` is connected to `/mnt/drive1` using mergerfs, if files exists in `/mnt/cache` it is prefered over `/mnt/drive1`

### Mergerfs setup

**mergerfs command:**

`mergerfs -o defaults,allow_other,use_ino,cache.files=off,category.create=epmfs /mnt/cache=RO:/mnt/drive1 /mnt/media`

**config in `/etc/fstab`:**

`/mnt/cache=RO:/mnt/drive1 /mnt/media fuse.mergerfs defaults,allow_other,use_ino,cache.files=off,category.create=epmfs 0 0`

### Database

Using a redis store with this conf `notify-keyspace-events Ex` so that it sends subscriber events for expiring keys.
