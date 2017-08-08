FROM golang:1.8

#make sure the neessary files and folders are available
RUN mkdir -p /var/log/nginx/
RUN touch /var/log/nginx/access.log
RUN touch /var/log/stats.log

#make $GOPATH/src/app the working area
RUN mkdir -p /go/src/nginxscraper
WORKDIR /go/src/nginxscraper
COPY nginxScraper/ nginxScraper/
COPY parsenginx/ parsenginx/

#main is located in nginxScraper folder
WORKDIR /go/src/nginxscraper/nginxScraper

RUN go build -o scraper .

CMD ["./scraper"]


