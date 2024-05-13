FROM debian:stable AS build-env

COPY config/tls/dod-wcf-root-ca-1.pem /usr/local/share/ca-certificates/dod-wcf-root-ca-1.pem.crt
COPY config/tls/dod-wcf-intermediate-ca-1.pem /usr/local/share/ca-certificates/dod-wcf-intermediate-ca-1.pem.crt
RUN apt-get update
# hadolint ignore=DL3008
RUN apt-get install -y ca-certificates --no-install-recommends
RUN update-ca-certificates

# hadolint ignore=DL3007
FROM gcr.io/distroless/base-debian11@sha256:2fb55308ef768a0ca0851f294d7f5b582579dba6522d1d2162e2d5f33b876e97

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY bin/rds-ca-rsa4096-g1.pem /bin/rds-ca-rsa4096-g1.pem
COPY bin/rds-ca-2019-root.pem /bin/rds-ca-2019-root.pem
COPY bin/milmove /bin/milmove

COPY config/tls/milmove-cert-bundle.p7b /config/tls/milmove-cert-bundle.p7b
COPY config/tls/dod-sw-ca-66.pem /config/tls/dod-sw-ca-66.pem

COPY swagger/* /swagger/
COPY build /build

ENTRYPOINT ["/bin/milmove"]

CMD ["serve", "--logging-level=debug"]

EXPOSE 8080
