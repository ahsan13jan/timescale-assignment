Feature: Benchmark CLI

  Scenario: Aggregate with single host name in one minute
    Given following data is in cpu usage table
      | host              | time                   | usage  |
      | host_000004       | 2017-10-12T07:00:00.00Z| 28.73  |
      | host_000004       | 2017-10-12T07:00:20.00Z| 58.73  |
      | host_000004       | 2017-10-12T07:00:40.00Z| 18.73  |
      | host_000004       | 2017-10-12T07:00:40.00Z| 18.73  |
    When cli tool run with following input csv path ../../testdata/query_params.csv
    And num of quries in the metrics should be 200


  Scenario: Aggregate with single host name in one minute
    Given following data is in cpu usage table
      | host              | time                   | usage  |
      | host_000004       | 2017-10-12T07:00:00.00Z| 28.73  |
      | host_000004       | 2017-10-12T07:00:20.00Z| 58.73  |
      | host_000004       | 2017-10-12T07:00:40.00Z| 18.73  |
      | host_000004       | 2017-10-12T07:00:40.00Z| 18.73  |
    When cli tool run with following input csv path ../../testdata/query_params_unit.csv
    And num of quries in the metrics should be 1
