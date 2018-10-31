FROM gcr.io/distroless/base

COPY bin/webserver /bin/mymove-server
COPY bin/chamber /bin/chamber

COPY config /config
COPY swagger/* /swagger/
COPY build /build

ENTRYPOINT ["/bin/mymove-server"]

CMD ["--debug-logging"]

EXPOSE 8080
