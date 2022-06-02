# Flink Data
[![Made by FilSwan](https://img.shields.io/badge/made%20by-FilSwan-green.svg)](https://www.filswan.com/)
[![Chat on Slack](https://img.shields.io/badge/slack-filswan.slack.com-green.svg)](https://filswan.slack.com)
[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg)](https://github.com/RichardLitt/standard-readme)

- Join us on our [public Slack channel](https://filswan.slack.com) for news, discussions, and status updates. 
- [Check out our medium](https://filswan.medium.com) for the latest posts and announcements.

Filnk data scans the data from Filecoin network and add important message to the database fro globle access.

## Table of Contents

- [Features](#Features)
- [Prerequisite](#Prerequisite)
- [Installation](#Installation)
- [Config](#Config)
- [License](#license)

## Features

Flink data provides the following functions:

* Get deal metadata from calibration 
* Convert units for some fields
* Store them into local database
* Provides web api to query deal metadata in our own format

## Prerequisite
- mysql database

## Installation
### Option:one: **Prebuilt package**

### Option:two: Source Code
:bell:**go 1.16+** is required
```shell
git clone https://github.com/filswan/flink.git
cd flink
git checkout <release_branch>
# create tables using scripts `./data/database/chain_link.sql` manually
cd data
./build_from_source.sh
```

### :bangbang: Important
After installation, flink-data maybe quit due to lack of configuration. Under this situation, you need
- :one: Edit config file `~/.swan/flink/data/config.toml` to solve this.
- :two: Execute `flink-data` using one of the following commands
```shell
./build/flink-data        #After installation from Option 2
```

### Note
- Logs are in directory ./logs
- You can add `nohup` before `./flink-data` to ignore the HUP (hangup) signal and therefore avoid stop when you log out.
- You can add `>> flink-data.log` in the command to let all the logs output to `flink-data.log`.
- You can add `&` at the end of the command to let the program run in background.
- You can only pass either calibration or mainnet as parameter
- Such as:
```shell
nohup ./build/flink-data calibration|mainnet >> flink-data.log &        #After installation from Option 2 (For calibration)       
```
- You can provide customized configuration file with -c flag
- Such as:
```shell
nohup ./build/flink-data calibration|mainnet -c [path_to_config] >> flink-data.log &        #After installation from Option 2 (For calibration)       
```

## Config
- **port**: Default `8886`, web api port for extension in future
- **release**: When work in release mode, set this to true, otherwise to false
### [main]
- **db_host**: Ip of the host for database instance running on
- **db_port**: Port of the host for database instance running on
- **db_schema_name**: Database schema name for flink data
- **db_username**: Username to access the database
- **db_password**: Password to access the database
- **db_args**: Other arguments to access database
- **db_max_idle_conn_num**: Maximum number of connections in the idle connection pool

### [chain_link]
- **bulk_insert_chainlink_limit**: When got more than this number of deals, than bulk insert them to db
- **bulk_insert_interval_milli_sec**: When deals in buffer exist(s), and time interval from last insert time to now is not less than this number, than bulk insert them to db
- **deal_id_interval_max**: Max interval between neighbour interval id

## Verification 
```shell
curl -X GET -H "content-type:application/json" "http://localhost:8886/network/filecoin_calibration"
```
```shell
curl -X POST -H "content-type:application/json" "http://localhost:8886/deal" --data '{ "deal_id":"87000", "network_name":"filecoin_calibration"}'
```

## License

[Apache](https://github.com/filswan/go-swan-provider/blob/main/LICENSE)

