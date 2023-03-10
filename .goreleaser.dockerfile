FROM scratch
COPY glue-deploy /glue-deploy
ENTRYPOINT [ "/glue-deploy" ]