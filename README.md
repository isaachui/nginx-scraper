# Nginx-scraper repo


## Objective
To create a program that reads local nginx log at `/var/log/nginx/access.log`. This will happen every 5 seconds, summarizing new entries with a stats-d compatible message. The summary will be appended to `/var/log/stats.log`

## How to test

*Assumes MiniKube is setup correctly with docker environment*

1. **Build docker containers**

  working directory: nginx-scraper repo root

  Build image for scraper. This is the program that watched the logfile.
  ```
  docker build -t scraper:1 .
  ```

  Build image for test-server. This is the test nginx server.
  ```
  docker build -t test-server:1 test-server/
  ```

2. **Deploy Kubernetes Pod**

  working directory: nginx-scraper repo root

  Run Kubectl to bring up Pod

  ```
  kubectl create -f k8s/scraper-deployment.yaml
  ```
3. **Check Scraper Kubectl Logs**

  To check the logs of the scraper-container, run the following:
  ```
  kubectl logs server-and-watcher scraper-container
  ```

  To tailf the logs add a `-f` flag after logs
  ```
  kubectl logs -f server-and-watcher scraper-container
  ```


4. **Test Webpage Status**


##
* Go program + Dockerfile
* Test nginx server + Dockerfile
* Kubernetes Configuration

Scraper.go contains the main function.

package parsenginx contains helper functions to parse nginx log file given a log format string.
