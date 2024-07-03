FROM gcr.io/distroless/static-debian12

WORKDIR app

COPY configs/default.yaml .
COPY proxy .

CMD ["/app/proxy", "-c", "default.yaml"]