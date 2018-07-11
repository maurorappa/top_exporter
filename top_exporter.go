package main
 
import (
                "fmt"
                "flag"
                "github.com/prometheus/client_golang/prometheus"
                "github.com/prometheus/client_golang/prometheus/promhttp"
                "github.com/shirou/gopsutil/process"
                "log"
                _ "net/http/pprof"
                "net/http"
                "os"
                "sort"
                "strings"
                "time"
)
 
type ProcInfo struct {
                Name  string
                Usage float64
                Pid   int32
}
 
type ByUsage []ProcInfo
 
func (a ByUsage) Len() int      { return len(a) }
func (a ByUsage) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByUsage) Less(i, j int) bool {
                return a[i].Usage > a[j].Usage
}
 
var ProcStat = prometheus.NewGaugeVec(prometheus.GaugeOpts{
                Name: "ProcessStat",
                Help: "Process name and cpu usage",
},
                []string{"cpustats"},
)
 
var (
                verbose        bool
                command_lenght int
                elements       int
                procinfos      []ProcInfo
                prev_processes []string
)
 
func init() {
                prometheus.MustRegister(ProcStat)
}
 
func top() {
                processes, _ := process.Processes()
                for _, p := range processes {
                                a, _ := p.CPUPercent()
                                n, _ := p.Cmdline()
                                p := p.Pid
                                if len(n) > command_lenght {
                                                n = n[0:command_lenght]
                                }
                                if len(n) == 0 {
                                                n = fmt.Sprintf("kernel(%d)",int(p) )
                                }
                                procinfos = append(procinfos, ProcInfo{n, a, p})
                }
               // extract the most CPU intensive           
                sort.Sort(ByUsage(procinfos))
        // clean the with the previos
                for _,v := range prev_processes {
          ProcStat.DeleteLabelValues(v)     
       }
 
 
                for _, p := range procinfos[:elements] {
                                if verbose {
                                                log.Printf("(%d) %s -> %f", int(p.Pid), p.Name, p.Usage)
                                }
                                ProcStat.WithLabelValues(p.Name).Set(float64(p.Usage))
                                prev_processes = append(prev_processes, p.Name)
                }
                procinfos = nil
}
 
func main() {
                flag.IntVar(&elements, "n", 10, "process to display")
                flag.IntVar(&command_lenght, "c", 50, "max character lenght of the command line")
                addr := flag.String("l", ":9097", "The address and port to listen on for HTTP requests.")
                interval := flag.Int("d", 5, "update interval")
                flag.BoolVar(&verbose, "v", false, "verbose on")
                flag.Parse()
                log.Printf("Metrics will be exposed on %s\n", *addr)
 
                if *interval > 0 {
                                //set a x seconds ticker
                                ticker := time.NewTicker(time.Duration(*interval) * time.Second)
 
                                go func() {
                                                for t := range ticker.C {
                                                                if verbose {
                                                                                log.Println("\nStats at", t)
                                                                }
                                                                top()
                                                }
                                }()
                }
 
                http.Handle("/metrics", promhttp.Handler())
                http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
                                w.Write([]byte(`
         <html>
         <head><title>Top Exporter</title></head>
         <body>
         <h1>Top Exporter</h1>
         <h2>parameters '` + strings.Join(os.Args, " ") + `'</h2>
         <p><a href='/metrics'><b>Metrics</b></a></p>
         </body>
         </html>
         `))
                })
                log.Fatal(http.ListenAndServe(*addr, nil))
 
}
