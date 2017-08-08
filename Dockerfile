FROM golang:1.8

#make sure the neessary files and folders are available
RUN mkdir -p /var/log/nginx/

RUN touch /var/log/nginx/access.log
RUN touch /var/log/stats.log

#forward writing to /var/log/stats.log to /dev/stdout
RUN ln -sf /dev/stdout /var/log/stats.log

#make $GOPATH/src/app the working area
RUN mkdir -p /go/src/nginxscraper
WORKDIR /go/src/nginxscraper
COPY scraper/ scraper/
COPY parsenginx/ parsenginx/

WORKDIR /go/src/nginxscraper/nginxScraper

RUN go-wrapper download
RUN go-wrapper install

##RUN go build -o /go/bin/scraper 

CMD ["go-wrapper", "run"]



