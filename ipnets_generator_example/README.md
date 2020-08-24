# IPNets Generator Example

This program uses the `ipnetsgenerator` module to generate a random list of IP nets.

## How to compile

```
$ make
```

## Execution examples

- Generate 10 random /32 IP nets

```
$ ./example -count=10
89.172.52.179/32
156.22.155.54/32
134.192.235.252/32
120.106.240.213/32
4.126.232.44/32
15.83.101.200/32
47.5.197.181/32
2.9.63.156/32
122.76.138.107/32
26.214.48.25/32
```

- Generate 10 random IP nets, 30% /24, 10% /16 and the remaining /32.

```
./example -count=10 -ipnets=24:0.3,16:0.1
11.104.4.0/24
174.20.130.0/24
163.81.80.0/24
23.133.0.0/16
26.6.193.29/32
79.160.36.132/32
134.134.35.73/32
102.22.219.30/32
222.96.102.224/32
72.48.253.38/3
```
