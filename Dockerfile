FROM harbor.csde.caci.com/docker.io/debian:stable AS build-env

COPY config/tls/dod-wcf-root-ca-1.pem /usr/local/share/ca-certificates/dod-wcf-root-ca-1.pem.crt
COPY config/tls/dod-wcf-intermediate-ca-1.pem /usr/local/share/ca-certificates/dod-wcf-intermediate-ca-1.pem.crt
RUN apt-get update
# hadolint ignore=DL3008
RUN apt-get install -y ca-certificates --no-install-recommends
RUN update-ca-certificates

# hadolint ignore=DL3007
FROM gcr.io/distroless/base-debian12@sha256:ad04bf079b9ed668d38fe2138cfe575847795985097b38a400f4ef1ff69a561a
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY bin/rds-ca-rsa4096-g1.pem /bin/rds-ca-rsa4096-g1.pem
COPY bin/rds-ca-2019-root.pem /bin/rds-ca-2019-root.pem
COPY bin/milmove /bin/milmove

COPY config/tls/milmove-cert-bundle.p7b /config/tls/milmove-cert-bundle.p7b
COPY config/tls/dod-sw-ca-75.pem /config/tls/dod-sw-ca-75.pem
COPY config/tls/dod-sw-ca-66.pem /config/tls/dod-sw-ca-66.pem

COPY swagger/* /swagger/
COPY build /build
COPY public/static/react-file-viewer /public/static/react-file-viewer

# Mount mutable tmp for app packages like pdfcpu
# hadolint ignore=DL3007
VOLUME ["/tmp"]

ENTRYPOINT ["/bin/milmove"]

CMD ["serve", "--logging-level=debug"]

EXPOSE 8080
