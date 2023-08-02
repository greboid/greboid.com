---
title: Looking for a new reverse proxy
date: 2022-03-23T19:55:02Z
---
The majority of container setups need a reverse proxy to redirect incoming requests to the appropriate container. I use [Dotege](https://github.com/csmith/dotege), a template generator that listens to container events and creates a config, this gives me the flexibility to switch proxies easily and try new things.

I was using Haproxy, but hit some irritating issues after a few releases and had to keep patching it to continue working.  One bug prevented me from pushing containers to my registry, and another caused issues with my site's background rendering. Both changes were caused some major changes to how they handled one thing or another, but there are no tests to make sure it works correctly.  Given this, I thought it was time to evaluate a replacement.

<!--more-->

## Benchmarking
To find the best proxy server, I decided to spin up a few instances of Dotege with different templates and run some benchmarks. I wanted to see how each proxy handled running my site and how easy they were to configure. I ran a 10-second test with 10 threads using [wrk](https://github.com/wg/wrk), a benchmarking tool.  Whilst performance isn't everything, it is a fairly important part of the job.  Here are the results of the benchmark tests:

### [Haproxy](https://www.haproxy.org/)
|                | Avg     | Stdev  | Max      | +/- Stdev |
|----------------|---------|--------|----------|-----------|
| Latency        | 19.78ms | 8.64ms | 139.07ms | 75.39%    |
| Thread Req/sec | 50.95   | 12.40  | 180.00   | 82.14%    |
| Total Req/sec  | 5193.85 |        |          |           |

{% link './haproxy.txt', 'Config' %}
&nbsp;
### [Centauri](https://github.com/csmith/centauri)
|                | Avg     | Stdev   | Max     | +/- Stdev |
|----------------|---------|---------|---------|-----------|
| Latency        | 24.52ms | 11.16ms | 89.78ms | 70.00%    |
| Thread Req/sec | 40.99   | 10.62   | 120.00  | 63.67%    |
| Total Req/sec  | 4163.63 |         |         |           |

{% link './centauri.txt', 'Config' %}
&nbsp;
### [Apache](https://httpd.apache.org/)
|                | Avg     | Stdev   | Max      | +/- Stdev |
|----------------|---------|---------|----------|-----------|
| Latency        | 25.12ms | 13.36ms | 190.05ms | 72.82%    |
| Thread Req/sec | 40.29   | 13.93   | 140.00   | 78.02%    |
| Total Req/sec  | 4097.49 |         |          |           |

{% link './apache.txt', 'Config' %}
&nbsp;
### [Nginx](https://nginx.org)
|                | Avg     | Stdev   | Max      | +/- Stdev |
|----------------|---------|---------|----------|-----------|
| Latency        | 29.65ms | 11.71ms | 198.40ms | 73.94%    |
| Thread Req/sec | 33.75   | 9.17    | 111.00   | 75.97%    |
| Total Req/sec  | 3412.71 |         |          |           |

{% link './nginx.txt', 'Config' %}
&nbsp;
### [Traefik](https://traefik.io/traefik/)
|                | Avg     | Stdev   | Max      | +/- Stdev |
|----------------|---------|---------|----------|-----------|
| Latency        | 42.83ms | 30.42ms | 455.54ms | 90.41%    |
| Thread Req/sec | 25.69   | 9.44    | 120.00   | 74.22%    |
| Total Req/sec  | 2529.13 |         |          |           |

{% link './traefik-static.txt', 'Static Config' %} {% link './traefik-dynamic.txt', 'Dynamic Config' %}
&nbsp;
### [Caddy](https://caddyserver.com/)
|                | Avg     | Stdev   | Max      | +/- St dev |
|----------------|---------|---------|----------|------------|
| Latency        | 50.25ms | 18.75ms | 159.01ms | 68.96%     |
| Thread Req/sec | 19.83   | 6.84    | 80.00    | 60.57%     |
| Total Req/sec  | 2017.68 |         |          |            |

{% link './caddy.txt', 'Config' %}
&nbsp;

## Summary
- Haproxy: As expected, Haproxy performed well. It's a well-established and optimized proxy server.  Haproxy isn't friendly to configure and it's very easy to introduce errors when you're trying new things.
- Apache and Nginx: Surprisingly, Nginx didn't perform as well as I expected. I spent some time optimizing it, but it still didn't match up to Apache. This made me re-evaluate my long-standing opinion that Nginx was faster than Apache.  Config of these is fairly easy, and there's a huge community and loads of examples.
- Caddy and Traefik: These two proxies had similar performance. They are both standard and modern Go proxies with similar goals. However, Caddy offers more features for web serving.  Caddy is very nice to configure and very flexible.  Traefik isn't fun to configure, once you've got your head round things it makes sense, but it always seems to be a struggle.
- Centauri: This came in a close second, which was a surprise after the Caddy and Traefik results, and just ahead of apache and nginx.  Configuration for this is very simple as it does one job, reverse proxy domains to another web server.

Looking at these results, it was a fairly easy choice to go with Centauri, it was the right combination of performance and ease to configure.  It also doesn't help it's written by the author of Dotege, [Chris Smith](https://chameth.com/), so getting support for bugs should be fairly straightforward.   This has been running a few weeks now and has yet to cause me any issues, so far so good.
