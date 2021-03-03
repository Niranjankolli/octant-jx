FROM scratch
EXPOSE 8080
ENTRYPOINT ["/octant-jx"]
COPY ./bin/ /