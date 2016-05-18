FROM scratch
COPY ./bin/blackflowhub /bin/zmarlin/
VOLUME ["/var/lib/blackflowhub"]
EXPOSE 5050
ENV BFH_BIND_ADDR=":5050" BFH_DB_CONN_STR="mongo:27017" BFH_FS_LOCATION="/var/lib/blackflowhub" BFH_JWT_SECRET=""
ENTRYPOINT ["/bin/zmarlin/blackflowhub"]