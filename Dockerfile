FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-demisto"]
COPY baton-demisto /