# gator

## Gator is an RSS aggreGATOR (Get it?) CLI that I wrote as part of a guided project on Boot.dev, a backend programming course.

### To use gator, you'll need to install Postgres and Go toolchain.

macOS with brew:

`brew install postgresql@15`

Linux / WSL (Debian). Here are the [docs from Microsoft](https://learn.microsoft.com/en-us/windows/wsl/tutorials/wsl-database#install-postgresql), but simply:

```
sudo apt update
sudo apt install postgresql postgresql-contrib
```

Ensure the installation worked. The `psql` command-line utility is the default client for Postgres. Use it to make sure you're on version 15+ of Postgres:
```psql --version```

(Linux only) Update postgres password:
```sudo passwd postgres```

Enter a password, and be sure you won't forget it. You can just use something easy like `postgres`.

### Next, install the Go toolchain:

**Option 1**: [The webi installer](https://webinstall.dev/golang/) is the simplest way for most people. Just run this in your terminal:

```bash
curl -sS https://webi.sh/golang | sh
```

_Read the output of the command and follow any instructions._

**Option 2**: Use the [official installation instructions](https://go.dev/doc/install).

Run `go version` on your command line to make sure the installation worked.

### Next, install the gator CLI in the root directory:

`go install gator`


### Gator is only configured to run locally. It needs a config file on your computer to work.


If there isn't one already, create a gatorconfig.json file in gator's root directory, like so:

```
{
    "db_url": "connection_string_goes_here",
    "current_user_name": "username_goes_here"
}
```

This keeps track of who is currently "logged" in, and the connection credentials for the Postgres database.

A connection string looks like this: protocol://username:password@host:port/database

Some examples:

macOS (no password, your username): `postgres://theotg:@localhost:5432/gator`
Linux (password set to postgres, user also set to postgres): `postgres://postgres:postgres@localhost:5432/gator`

You can test your connection string by running

`psql "postgres://wagslane:@localhost:5432/gator"`

Note: There's no user-based authentication for this app. If someone has the database credentials, they can act as any user. 

### Gator Commands:

Gator can take a number of commands, but you must first `register` a username and `login` to set that name as the active user, so you can follow feeds.

`register <username>` - register a new user.

`login <username>` - log in as the provided user.

`users` - retrieves a list of registered users in the database.

`addfeed <feed name> <url>` - add a feed to the database.

`agg <time duration string: 1h, 1m, 5s>` - begin aggregating the users feeds for browsing. CAUTION - do not set the refresh duration to be too short to prevent DOS-ing the server!

`browse <post limit>` - browse the user's aggregated post by a given number. Defaults to two posts per browse if not given the limit.

`feeds` - retrieves a list of feeds in the database.

`follow <url>` - follow a feed with the given URL.

`unfollow <url>` unfollow a feed at the given URL.

`following` - retrieves a list of the current user's followed feeds.

`reset` - resets the database.
