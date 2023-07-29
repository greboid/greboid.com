---
title: "Web Proxy"
date: 2022-03-23T19:55:02Z
resources:
- name: "apache.txt"
- name: "caddy.txt"
- name: "centauri.txt"
- name: "haproxy.txt"
- name: "nginx.txt"
- name: "traefik-static.txt"
- name: "traefik-dynamic.txt"
---
The majority of container setups need a reverse proxy to redirect incoming requests to the appropriate container. Some 
time ago I went with a template generator written by my friend, [Dotege](https://github.com/csmith/dotege), this avoids
having to expose the docker socket to the web and makes the choice of proxy reasonable flexible, but by default it ships
with a template for Haproxy - as a very established proxy, I went with this and this served me quite well for some time.

I have more recently experienced some annoying bugs with haproxy, one of which prevented any containers being pushed 
to my registry and a more recent one caused an intermittent issue with my site not rendering a background due to a 
permissions-policy header causing an internal server error.  In light of this I took the opportunity to take a look at 
the alternatives, and as everything is in docker and the configs can be templated as needed, I decided the best way to 
pick a proxy would be to spin up a few instances of Dotege with different templates and see how they handled running my
site with some simple benchmarks and pick one that was nicest to configure whilst being quite performant.

The below results were a 10-second run, with 10 threads using [wrk](https://github.com/wg/wrk), whilst the configs for
these might not be the most optimised examples, I think they're fairly optimised and definitely typical examples. 
This should give a look at both the latency and throughput of the proxy, I picked a few big examples and them some less 
common ones, the results were not exactly as expected.

### [Haproxy](https://www.haproxy.org/)
|                | Avg     | Stdev  | Max      | +/- Stdev |
|----------------|---------|--------|----------|-----------|
| Latency        | 19.78ms | 8.64ms | 139.07ms | 75.39%    |
| Thread Req/sec | 50.95   | 12.40  | 180.00   | 82.14%    |
| Total Req/sec  | 5193.85 |        |          |           |

[Config](haproxy.txt)
&nbsp;
### [Centauri](https://github.com/csmith/centauri)
|                | Avg     | Stdev   | Max     | +/- Stdev |
|----------------|---------|---------|---------|-----------|
| Latency        | 24.52ms | 11.16ms | 89.78ms | 70.00%    |
| Thread Req/sec | 40.99   | 10.62   | 120.00  | 63.67%    |
| Total Req/sec  | 4163.63 |         |         |           |

[Config](centauri.txt)
&nbsp;
### [Apache](https://httpd.apache.org/)
|                | Avg     | Stdev   | Max      | +/- Stdev |
|----------------|---------|---------|----------|-----------|
| Latency        | 25.12ms | 13.36ms | 190.05ms | 72.82%    |
| Thread Req/sec | 40.29   | 13.93   | 140.00   | 78.02%    |
| Total Req/sec  | 4097.49 |         |          |           |

[Config](apache.txt)
&nbsp;
### [Nginx](https://nginx.org)
|                | Avg     | Stdev   | Max      | +/- Stdev |
|----------------|---------|---------|----------|-----------|
| Latency        | 29.65ms | 11.71ms | 198.40ms | 73.94%    |
| Thread Req/sec | 33.75   | 9.17    | 111.00   | 75.97%    |
| Total Req/sec  | 3412.71 |         |          |           |

[Config](nginx.txt)
&nbsp;
### [Traefik](https://traefik.io/traefik/)
|                | Avg     | Stdev   | Max      | +/- Stdev |
|----------------|---------|---------|----------|-----------|
| Latency        | 42.83ms | 30.42ms | 455.54ms | 90.41%    |
| Thread Req/sec | 25.69   | 9.44    | 120.00   | 74.22%    |
| Total Req/sec  | 2529.13 |         |          |           |

[Static Config](traefik-static.txt) [Dynamic Config](traefik-dynamic.txt)
&nbsp;
### [Caddy](https://caddyserver.com/)
|                | Avg     | Stdev   | Max      | +/- St dev |
|----------------|---------|---------|----------|------------|
| Latency        | 50.25ms | 18.75ms | 159.01ms | 68.96%     |
| Thread Req/sec | 19.83   | 6.84    | 80.00    | 60.57%     |
| Total Req/sec  | 2017.68 |         |          |            |

[Config](caddy.txt)
&nbsp;

Haproxy is an obvious leader here, it's a very well optimised long-established proxy, with no intentions to be anything 
else, so seeing it come out on top was no surprise.

Apache and nginx definitely did not come out in the order I expected, I did spend some time optimising nginx as it 
seemed like an outlier, I did manage to increase its performance by about 50% from my initial config.  I have for a 
long time held the opinion that nginx was a faster alternative to apache when it came to being a webserver and proxy, 
but after seeing these I have definite re-evaluated these opinions.

Caddy and Traefik seem to be about the same, which makes sense as they're both fairly standard and modern Go proxies 
with almost the same goals - although caddy does do an awful lot more than Traefik in terms of web serving.

I had initially started using Caddy, its defaults are very sensible and config file is terse yet flexible, but after
running these benchmarks have since made a switch to Centauri (and switched my website from Caddy to Apache!).  This is
a very new proxy server (written by [Chris Smith](https://chameth.com/) the author of Dotege), its considerably more 
performant then the others here and still meets all my requirements and (touch wood) isn't showing any hard to 
reproduce bugs in what should be fairly normal web traffic.
