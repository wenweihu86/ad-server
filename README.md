# ad-server
advertising server
# the third lib
- github.com/ibbd-dev/go-async-log
- github.com/satori/go.uuid
- github.com/spf13/viper
# benchmark
## Hardware
CPU: 12 cores<br>
Memory: 64G
## Results:
Used Connections:               1000<br>
Used Threads:                   12<br>
Total number of calls:          1000000<br>

===========================TIMINGS===========================<br>
Total time passed:              46.59s<br>
Avg time per request:           44.09ms<br>
Requests per second:            21464.89<br>
Median time per request:        42.80ms<br>
99th percentile time:           63.66ms<br>
Slowest time for request:       3002.00ms<br>

=============================DATA=============================<br>
Total response body sizes:              1202000594<br>
Avg response body per request:          1202.00ms<br>
Transfer rate per second:               25800812.03 Byte/s (25.80 MByte/s)<br>
==========================RESPONSES==========================<br>
20X Responses:          999998  (100.00%)<br>
30X Responses:          0       (0.00%)<br>
40X Responses:          0       (0.00%)<br>
50X Responses:          0       (0.00%)<br>
Errors:                 2       (0.00%)<br>
