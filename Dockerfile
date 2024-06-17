FROM golang:1.22 as build
WORKDIR /code
COPY . .
RUN make build
RUN chmod +x matchmaker
CMD ["./matchmaker"]