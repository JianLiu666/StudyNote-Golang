# NATS Benchmark Report

```shell
# Case1: Run a 1:N throughput test
nats bench foo --pub 1 --sub 5 --size 16 --msgs 10000000
```

| version | min (msgs/sec) | avg (msgs/sec) | max (msgs/sec) |
| ------- | -------------: | -------------: | -------------: |
| 2.9.3   | 899,498        | 899,523        | 899,540        |
| 2.9.16  | 848,201        | 848,265        | 848,318        |

```shell
# Case2: Run a 1:N throughput test (JetStream)
nats bench foo --js --pub 1 --sub 5 --size 16 --msgs 1000000
```

| version | min (msgs/sec) | avg (msgs/sec) | max (msgs/sec) |
| ------- | -------------: | -------------: | -------------: |
| 2.9.3   | 68,543         | 68,591         | 68,626         |
| 2.9.16  | 83,283         | 83,332         | 83,370         |