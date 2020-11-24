FROM ubuntu:16.04
LABEL maintainer="Prashant Iyer"

WORKDIR /

COPY bin/dispatch /bin/
COPY dispatch.yaml /

ENTRYPOINT ["dispatch", "run", "dispatch.yaml"]