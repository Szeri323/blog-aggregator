# Aggregator Feeds project

## Short description
The [Boot.dev](https://www.boot.dev/courses/build-blog-aggregator-golang) project provaides set of skill like using postgres, orm like goose and sqlc, using middleware, and networking.

## What you need
To use this program you need to have had installed Go and PostgreSQL on your machine.
To install Go:
checkout official go installation page https://go.dev/doc/install
```bash 
sudo apt install go
``` 
to check if it is installed you can type:
```bash
go version
```
and to install postgres:
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```
to check if it is installed you can type:
```bash
psql --version
```

## How to build/install
You need to bo in gator root directory. Next you type:
```bash
go build
```
To create a binary file and use it in root directory via relative:
```bash 
./gator <...rest of commands>
```
or absolute path: 
```bash 
~/gator <...rest of commands>
```
To install gator CLI app in to your system and use it where every path you want you need to execute go install command in root gator directory:
```bash
go install
```

## Basic instruction
To program works popreply you need to create .gatorconfig.json file.
To use gator command you need to follow the following template:
```bash 
gator command [...options]
```
Basic command that you can execute are:
- login [user_name:string] - switch user,
- register [user_name:string] - register new user and set it as active,
- users - list all existing users with mark avtive one,
- reset - (only for debbuging and if you know what you are doing),
- addfeed [feed_name, feed_url] - add feed and start following it by current user,
- feeds - list all feeds,
- follow [feed_url] - follow the feed,
- following - list feeds that are followed by current user,
- unfollow [feed_url] - unfollow the feed by current user,
- agg - fetch feed content from feeds to posts table in inifinite loop (to terminate loop use ctrl + c combination),
- browse [limit:int] - List posts for current user with limit,

