# Pod scraper

Pod scraper is an application that scrapes the pods discovered in the namespaces labeled with the value of -ns-label and/or pods labeled with the value of -pod-label.
This small application is meant to run from a pod as it scrapes the pod's internal IP and loads the client-go in-cluster configuration.

```
$ ./pod-scraper -help
Usage of ./pod-scraper:
  -code int
        Expected status code (default 200)
  -endpoint string
        Target endpoint (default "/")
  -ns-label string
        Target namespace label
  -pod-label string
        Target pod label
  -port int
        Target port (default 80)
  -scheme string
        URL scheme, http or https (default "http")
  -timeout duration
        Request timeout (default 10s)
```
