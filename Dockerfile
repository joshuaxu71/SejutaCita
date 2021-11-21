FROM golang:1.17.3
RUN mkdir /app
ADD . /app
WORKDIR /app
ENV PORT=:9090
ENV DEPLOYMENT=development
ENV DB_NAME=SejutaCita
ENV DEV_DB=mongodb://192.168.49.2:32000
ENV SECRET_KEY=oskgkffugka
RUN go build -o main .
CMD ["/app/main"]