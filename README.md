# Parquet to CSV

Have a parquet file? Want a CSV? Got you covered.

## Installation (for Mac OS)

1. Download [the latest release for your platform (Darwin)](https://github.com/eliblock/parquet-to-csv/releases/latest).
1. Double-click the downloaded file to un-tar it. A folder will open.
1. In the folder, right (two-finger) click on `parquet-to-csv` and then select `Open` from the context menu.
1. In the pop-ups, elect to open the file _anyways_. This instructs Mac OS Gatekeeper to not block the executable even though it is not signed by a known developer.
1. _Either_ move the `parquet-to-csv` file somewhere on your `PATH`, _or_ place it somewhere memorable on your machine (e.g., `~/Desktop/parquet-to-csv`).

## Usage

```sh
parquet-to-csv <path_to_file>
```

For example:

```sh
parquet-to-csv ./data.parquet
```

Or, if you placed `parquet-to-csv` somewhere not on your `PATH`, use that filepath. For example, if you placed it at `~/Desktop/parquet-to-csv`...

```sh
~/Desktop/parquet-to-csv <path_to_file>
```

## Configuration

### `-n`

The `-n` flag sets a maximum number of rows to be processed.

Defaults to all rows

```sh
# Outputs at most 5 rows
parquet-to-csv -n 5 ./data.parquet
```

### `--in`

The `--in` flag sets the file path where the parquet file may be read from. It is mutually exclusive with using an argument, but may be used when expressive CLI args are desired.

```sh
# These both process the file ./data.parquet
parquet-to-csv ./data.parquet
parquet-to-csv --in ./data.parquet
```

### `--out`

The `--out` flag sets the file path where the output csv should be written. Must be a write-able file path.

Defaults to standard output

```sh
# Outputs to /tmp/out.csv
parquet-to-csv --out /tmp/out.csv ./data.parquet
```

#### Alternative

`parquet-to-csv`, by default, outputs csv rows to standard output (alongside log messages to standard error). Instead of specifying a `--out` file, you may redirect the output.

```sh
# These both write output to /tmp/out.csv
parquet-to-csv --out /tmp/out.csv ./data.parquet
parquet-to-csv ./data.parquet > /tmp/out.csv
```

Relatedly, log messages may be suppressed by suppressing standard error.

```sh
parquet-to-csv ./data.parquet 2> /dev/null
```

### `--overwrite`

By default, `parquet-to-csv` prevents use of a file specified by the `--out` flag which already exists. If you would prefer to overwrite a file if it already exists, specify `--overwrite`.

Defaults to false (no overwrite permitted)

```sh
cat "hello, world" > /tmp/out.csv
# Overwrites /tmp/out.csv
parquet-to-csv --out /tmp/out.csv --overwrite ./data.parquet
```

## Development

### Environment

```sh
brew install go@1.19
go build ./...
go test -v ./...
```

### Release

```sh
brew install goreleaser
git tag v0.1.0 # update for your version
git push origin v0.1.0 # update for your version
goreleaser release --rm-dist --snapshot # remove --snapshot for a full release
# complete the release on GitHub
```
