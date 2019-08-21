Use symlinks to divide files from one directory into multiple directories, and recombine as well.

# Usage

Split up a directory

```
% fdivide --help
fdivide
Divide a regular files from a directory into subdirectories by number of files using symlinks.

Usage:
    fdivide --size <dir-size> <input-dir> <output-dir>
    fdivide --into <num-dirs> <input-dir> <output-dir>
```

And put it back together

```
% fcombine --help
fcombine
Combine files from sibiling subdirectories into a single output directory using symlinks.

Usage:
    fcombine <input-parent-dir> <output-dir>
```

# Examples

Create some empty files

```
% mkdir initial
% for f in {1..20}; do touch initial/$(echo $f | xargs printf '%02d'); done
% ls initial
01  02  03  04  05  06  07  08  09  10  11  12  13  14  15  16  17  18  19  20
```

Split them into tree directories
```
% fdivide --into 3 initial divided-into
% ls divided-into -lR
divided-into:
total 12
drwxr-xr-x 2 oliver oliver 4096 Aug 20 23:28 0
drwxr-xr-x 2 oliver oliver 4096 Aug 20 23:28 1
drwxr-xr-x 2 oliver oliver 4096 Aug 20 23:28 2

divided-into/0:
total 0
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 01 -> /home/oliver/files/initial/01
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 02 -> /home/oliver/files/initial/02
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 03 -> /home/oliver/files/initial/03
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 04 -> /home/oliver/files/initial/04
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 05 -> /home/oliver/files/initial/05
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 06 -> /home/oliver/files/initial/06
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 07 -> /home/oliver/files/initial/07

divided-into/1:
total 0
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 08 -> /home/oliver/files/initial/08
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 09 -> /home/oliver/files/initial/09
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 10 -> /home/oliver/files/initial/10
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 11 -> /home/oliver/files/initial/11
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 12 -> /home/oliver/files/initial/12
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 13 -> /home/oliver/files/initial/13
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 14 -> /home/oliver/files/initial/14

divided-into/2:
total 0
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 15 -> /home/oliver/files/initial/15
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 16 -> /home/oliver/files/initial/16
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 17 -> /home/oliver/files/initial/17
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 18 -> /home/oliver/files/initial/18
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 19 -> /home/oliver/files/initial/19
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:28 20 -> /home/oliver/files/initial/20
```

Split them into directories which each contain four files
```
% fdivide --size 4 initial divided-size
% ls divided-size -lR
divided-size:
total 20
drwxr-xr-x 2 oliver oliver 4096 Aug 20 23:29 0
drwxr-xr-x 2 oliver oliver 4096 Aug 20 23:29 1
drwxr-xr-x 2 oliver oliver 4096 Aug 20 23:29 2
drwxr-xr-x 2 oliver oliver 4096 Aug 20 23:29 3
drwxr-xr-x 2 oliver oliver 4096 Aug 20 23:29 4

divided-size/0:
total 0
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 01 -> /home/oliver/files/initial/01
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 02 -> /home/oliver/files/initial/02
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 03 -> /home/oliver/files/initial/03
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 04 -> /home/oliver/files/initial/04

divided-size/1:
total 0
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 05 -> /home/oliver/files/initial/05
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 06 -> /home/oliver/files/initial/06
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 07 -> /home/oliver/files/initial/07
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 08 -> /home/oliver/files/initial/08

divided-size/2:
total 0
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 09 -> /home/oliver/files/initial/09
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 10 -> /home/oliver/files/initial/10
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 11 -> /home/oliver/files/initial/11
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 12 -> /home/oliver/files/initial/12

divided-size/3:
total 0
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 13 -> /home/oliver/files/initial/13
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 14 -> /home/oliver/files/initial/14
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 15 -> /home/oliver/files/initial/15
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 16 -> /home/oliver/files/initial/16

divided-size/4:
total 0
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 17 -> /home/oliver/files/initial/17
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 18 -> /home/oliver/files/initial/18
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 19 -> /home/oliver/files/initial/19
lrwxrwxrwx 1 oliver oliver 29 Aug 20 23:29 20 -> /home/oliver/files/initial/20
```

Recombine!

```
% fcombine divided-into combined-1
% ls -l combined-1
total 0
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 01 -> /home/oliver/files/divided-into/0/01
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 02 -> /home/oliver/files/divided-into/0/02
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 03 -> /home/oliver/files/divided-into/0/03
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 04 -> /home/oliver/files/divided-into/0/04
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 05 -> /home/oliver/files/divided-into/0/05
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 06 -> /home/oliver/files/divided-into/0/06
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 07 -> /home/oliver/files/divided-into/0/07
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 08 -> /home/oliver/files/divided-into/1/08
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 09 -> /home/oliver/files/divided-into/1/09
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 10 -> /home/oliver/files/divided-into/1/10
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 11 -> /home/oliver/files/divided-into/1/11
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 12 -> /home/oliver/files/divided-into/1/12
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 13 -> /home/oliver/files/divided-into/1/13
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 14 -> /home/oliver/files/divided-into/1/14
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 15 -> /home/oliver/files/divided-into/2/15
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 16 -> /home/oliver/files/divided-into/2/16
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 17 -> /home/oliver/files/divided-into/2/17
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 18 -> /home/oliver/files/divided-into/2/18
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 19 -> /home/oliver/files/divided-into/2/19
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 20 -> /home/oliver/files/divided-into/2/20
```

Recombine some more!

```
% fcombine divided-size combined-2
% ls -l combined-2
total 0
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 01 -> /home/oliver/files/divided-size/0/01
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 02 -> /home/oliver/files/divided-size/0/02
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 03 -> /home/oliver/files/divided-size/0/03
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 04 -> /home/oliver/files/divided-size/0/04
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 05 -> /home/oliver/files/divided-size/1/05
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 06 -> /home/oliver/files/divided-size/1/06
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 07 -> /home/oliver/files/divided-size/1/07
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 08 -> /home/oliver/files/divided-size/1/08
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 09 -> /home/oliver/files/divided-size/2/09
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 10 -> /home/oliver/files/divided-size/2/10
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 11 -> /home/oliver/files/divided-size/2/11
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 12 -> /home/oliver/files/divided-size/2/12
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 13 -> /home/oliver/files/divided-size/3/13
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 14 -> /home/oliver/files/divided-size/3/14
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 15 -> /home/oliver/files/divided-size/3/15
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 16 -> /home/oliver/files/divided-size/3/16
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 17 -> /home/oliver/files/divided-size/4/17
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 18 -> /home/oliver/files/divided-size/4/18
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 19 -> /home/oliver/files/divided-size/4/19
lrwxrwxrwx 1 oliver oliver 36 Aug 20 23:30 20 -> /home/oliver/files/divided-size/4/20
```
