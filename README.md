# ad-server
advertising server
# Getting Start
- install go and glide env
- cp -f mirrors.yaml ~/.glide/
- glide install
- go build -o ad_server main/main.go
- ./ad_server
# Benchmark
#### Hardware
CPU: 12 cores<br>
Memory: 64G
#### Benchmark Input
./hey -n 1000000 -c 2000 "http://localhost:8001/ad/search?slot_id=2&ad_num=1&ip=101.88.50.181&device_id=0x22q53&os=1&os_version=1.0.0" 
#### Benchmark Result:
- Total:        24.3260 secs
- Slowest:      0.7821 secs
- Fastest:      0.0002 secs
- Average:      0.0467 secs
- Requests/sec: 41108.3054
- Total data:   1201899425 bytes
- Size/request: 1201 bytes
