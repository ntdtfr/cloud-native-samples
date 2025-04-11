# syntax=docker/dockerfile:1

ARG KEYCLOAK_VERSION=26.1.4

###
### Stage 1: Keycloak base image
FROM quay.io/keycloak/keycloak:${KEYCLOAK_VERSION} AS build


# # for demonstration purposes only, please make sure to use proper certificates in production instead
# RUN keytool -genkeypair -storepass password -storetype PKCS12 -keyalg RSA -keysize 2048 -dname "CN=server" -alias server -ext "SAN:c=DNS:localhost,IP:127.0.0.1" -keystore conf/server.keystore

RUN /opt/keycloak/bin/kc.sh build

###
### Stage 3: Keycloak base image
FROM quay.io/keycloak/keycloak:${KEYCLOAK_VERSION} AS keycloak-base

# Enable health and metrics support
# https://www.keycloak.org/server/containers
ENV KC_HEALTH_ENABLED=true
ENV KC_METRICS_ENABLED=true

# Set working directory
WORKDIR /opt/keycloak

# Copy Keycloak from build stage to the image
COPY --chown=keycloak:keycloak --from=build /opt/keycloak/ .

###
### Stage 4: Keycloak development image
FROM keycloak-base AS development

# Set working directory
WORKDIR /opt/keycloak

# Copy local data to the image
COPY --chown=keycloak:keycloak ./data ./data

# Set entrypoint
ENTRYPOINT ["/opt/keycloak/bin/kc.sh", "start-dev"]
# ENTRYPOINT ["/opt/keycloak/bin/kc.sh", "start-dev", "--cache=local", "--import-realm"]
