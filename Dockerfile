FROM gcr.io/distroless/base:latest

COPY bin/rds-ca-2019-root.pem /bin/rds-ca-2019-root.pem
COPY bin/milmove /bin/milmove

COPY config/tls/Certificates_PKCS7_v5.6_DoD.der.p7b /config/tls/Certificates_PKCS7_v5.6_DoD.der.p7b
COPY config/tls/dod-sw-ca-54.pem /config/tls/dod-sw-ca-54.pem
COPY config/tls/devlocal-ca.key /config/tls/devlocal-ca.key
COPY config/tls/devlocal-ca.pem /config/tls/devlocal-ca.pem

COPY swagger/* /swagger/
COPY build /build

ENTRYPOINT ["/bin/milmove"]

CMD ["serve", "--debug-logging"]

EXPOSE 8080
