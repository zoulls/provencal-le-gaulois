FROM alpine:3.24

COPY bin/provencal-le-gaulois .
COPY .env-prod .env

ENTRYPOINT ["./provencal-le-gaulois"]
