FROM gcr.io/distroless/base:latest

COPY bin/rds-combined-ca-bundle.pem /bin/rds-combined-ca-bundle.pem
COPY bin/chamber /bin/chamber
COPY bin/milmove /bin/milmove

COPY config/tls/Certificates_PKCS7_v5.4_DoD.der.p7b /config/tls/Certificates_PKCS7_v5.4_DoD.der.p7b
COPY config/tls/dod-sw-ca-54.pem /config/tls/dod-sw-ca-54.pem

COPY swagger/* /swagger/
COPY build /build

ENTRYPOINT ["/bin/milmove"]

CMD ["serve", "--debug-logging"]

EXPOSE 8080
