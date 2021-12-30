FROM golang:latest
LABEL version="1.0" \
      "site.name"="FORUM" \
      maintainers="Rik0 and ALika96" \
      release-date="November 18, 2021" \
      description="FORUM" \
      authors="Rik0 and Alika96"
WORKDIR /forum

RUN mkdir model
COPY model model/

RUN mkdir pkg
COPY pkg pkg/

RUN mkdir routes
COPY routes routes/

RUN mkdir sqlite
COPY sqlite sqlite/

RUN mkdir static
COPY static static/

RUN mkdir views
COPY views views/

COPY config.json .
COPY go.mod .
COPY go.sum .
COPY main.go . 

RUN go build -o main .
CMD ["./main"]