FROM gcr.io/distroless/base:latest

COPY bin/rds-combined-ca-bundle.pem /bin/rds-combined-ca-bundle.pem

COPY bin/chamber /bin/chamber

COPY bin/milmove /bin/milmove

COPY config /config

COPY swagger/* /swagger/

COPY build /build

ENTRYPOINT ["/bin/milmove"]

EXPOSE 8080
