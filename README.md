# GeoIP Updater
*GeoIP Updater Microservice*

## Building
To build this project, you need my custom Docker build extension dockpipe, which you can get here [here](https://github.com/lorenz/dockpipe).
Then you can just type
```
$ dockpipe geoip-updater:dev .
```
to get an image built.

## Environment variables
| Name | Default | Description |
| ---- | ------- | ----------- |
| `USER_ID`| 99999 | The user id to use to connect to the MaxMind update server, the default one is anonymous for GeoLite databases |
| `EDITION_IDS` | GeoLite2-City | A comma-separated list of edition ids to download and keep up-to-date |

## Volumes
| Path | Description |
| ---- | ----------- |
| `/data` | Shared volume with geoip-server instances for storing GeoIP-Databases. Should be on SSD or high-IOPS volume. |

## Facts
* Only updates databases when they actually change to avoid rate limits and consume excessive bandwidth.
* Checks for updates every hour.
* Does replace databases atomically
 