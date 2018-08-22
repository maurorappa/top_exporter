# top_exporter
Prometheus exporter for the most CPU intensive Linux processes (top like)

The first version used github.com/shirou/gopsutil/process which is widely used by lots of monitoring tools, but I find out it's too CPU intensive (see also https://github.com/shirou/gopsutil/issues/413)
So I wrote another (top_exporter_agnostic) which scan /proc and list the processes based on CPU time used without any percentage calculation which varies wildly from tool to tool.

For reference see:
* https://rosettacode.org/wiki/Linux_CPU_utilization
* https://jaroslawr.com/articles/mastering-linux-performance-cpu-time-and-cpu-usage/
* https://github.com/uber-common/cpustat

