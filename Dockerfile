FROM alpine:3.14

COPY bin/provencal-le-gaulois .
COPY config/config.json .
COPY .env-prod .env

# For dev
#COPY .env .

ENTRYPOINT ["./provencal-le-gaulois"]
