FROM postgres:15

WORKDIR /src
COPY entrypoint.sh /entrypoint.sh
RUN chmod u+x /entrypoint.sh
RUN /entrypoint.sh