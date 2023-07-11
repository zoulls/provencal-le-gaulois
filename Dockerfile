FROM alpine:3.14

COPY bin/provencal-le-gaulois .
COPY .env-prod .env

ENTRYPOINT ["./provencal-le-gaulois"]
