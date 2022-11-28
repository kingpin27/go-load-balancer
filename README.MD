# Layer 4 TCP Load-balancer written in GO

A simple layer 4 TCP load-balancer written purely in GO using standard net package.

For the load balancing algorithm, a random number is generated less than number of servers available. If the server is not alive, next nearest alive server is selected. You can find more about algo in [balancingAlgo.go](/balancingAlgo.go).

# Benchmark

Benchmark is done on my peronal laptop(i5, 4cores, 8GB). All loadbalancer, 4 servers and tester(wrk) are running on same machine. wrk is running on wsl.

![wrk benchmark](/benchmark.png)