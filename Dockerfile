FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-xsoar"]
COPY baton-xsoar /