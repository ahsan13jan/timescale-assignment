Feature: Postgres Client
  Give the cpu usage the postgres client should aggregate the query

  Scenario: Aggregate with single host name in one minute
    Given following data is in cpu usage table
      | host              | time                   | usage  |
      | host_000004       | 2017-10-12T07:00:00.00Z| 28.73  |
      | host_000004       | 2017-10-12T07:00:20.00Z| 58.73  |
      | host_000004       | 2017-10-12T07:00:40.00Z| 18.73  |
    When aggregated per minute of host host_000004 result should be from 2017-10-12T07:00:00.00Z to 2017-10-12T07:01:00.00Z
     | time                   | min usage  | max usage  |
     | 2017-10-12T07:00:00.00Z| 18.73      | 58.73      |


  Scenario: Aggregate with single host name in two minute
    Given following data is in cpu usage table
      | host              | time                   | usage  |
      | host_000004       | 2017-10-12T07:00:00.00Z| 28.73  |
      | host_000004       | 2017-10-12T07:00:20.00Z| 58.73  |
      | host_000004       | 2017-10-12T07:01:20.00Z| 18.73  |
      | host_000004       | 2017-10-12T07:01:40.00Z| 8.73   |

    When aggregated per minute of host host_000004 result should be from 2017-10-12T07:00:00.00Z to 2017-10-12T07:02:00.00Z
      | time                   | min usage  | max usage  |
      | 2017-10-12T07:00:00.00Z| 18.73      | 58.73      |
      | 2017-10-12T07:01:00.00Z| 8.73       | 18.73      |


  Scenario: Aggregate with two host name in minute
    Given following data is in cpu usage table
      | host              | time                   | usage  |
      | host_000004       | 2017-10-12T07:00:00.00Z| 28.73  |
      | host_000004       | 2017-10-12T07:00:20.00Z| 58.73  |
      | host_000005       | 2017-10-12T07:01:20.00Z| 18.73  |
      | host_000005       | 2017-10-12T07:01:40.00Z| 8.73   |

    When aggregated per minute of host host_000004 result should be from 2017-10-12T07:00:00.00Z to 2017-10-12T07:02:00.00Z
      | time                   | min usage  | max usage  |
      | 2017-10-12T07:00:00.00Z| 18.73      | 58.73      |
      | 2017-10-12T07:01:00.00Z| 8.73       | 18.73      |
