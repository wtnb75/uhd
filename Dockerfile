FROM scratch
ARG EXT=
COPY uhd${EXT} /uhd
ENTRYPOINT ["/uhd"]
