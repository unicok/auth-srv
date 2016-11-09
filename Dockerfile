FROM alpine:3.2
ADD auth /auth
ENTRYPOINT [ "/auth" ]