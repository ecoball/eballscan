Eballscan
-------

## Depends
You need install [CockroachDB](https://www.cockroachlabs.com/docs/stable/install-cockroachdb.html) 


## Build
Run './eballscan_build.sh' in eballscan

If you have already executed this script,You can choose whether to delete old  eballscan and cockroach-data or not.

```bash
$:~/go/src/github.com/ecoball/eballscan$ ./eballscan_build.sh
```

## run
You have to start a full node of ecoball, and then start eballscan

The database is started during startup

```bash
$:~/go/src/github.com/ecoball/eballscan$ ./eballscan_service.sh
```