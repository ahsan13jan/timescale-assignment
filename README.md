# Timescale Assignment

# build cli
`make build `

## environment variables
- POSTGRES_DB
- POSTGRES_HOST
- POSTGRES_PORT
- POSTGRES_USERNAME
- POSTGRES_PASSWORD
- CSV_PATH
- MAX_CONCURRENT_WORKERS
- LOG_LEVEL

`CSV_PATH=path.csv POSTGRES_USERNAME=username ./benchmark`

## flags
- file path
- max workers

`./benchmark -f "./testdata/query_params.csv" -m 2`

## stdin
` cat testdata/query_params.csv | ./benchmark`


## debug cli
` LOG_LEVEL=debug ./benchmark -f "./testdata/query_params.csv" -m 2`

## BDD test and testing
BDD test can be found `internal/testing/features/`

Run all tests with
`make test-all`