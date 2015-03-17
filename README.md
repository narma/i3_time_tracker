# Worktime tracker for i3wm

## Setup and usage

```sh
go build -a
./i3_time_tracker 1 2 3 4
```
## How it works

Workspaces are divided to those which are used for work and those which aren't.
Time is tracking only when chosen workspaces are active.

Currently results synchronize with redis. It's up to you how to use them.
