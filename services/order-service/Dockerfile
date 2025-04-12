ARG NODE_VERSION=22

##
## Stage: Base image
FROM node:${NODE_VERSION}-alpine AS base

# Define environment variables
ENV APP_DIR=/opt/app
ENV NODE_PORT=3000

# Create user and group
# RUN apk --no-cache add ca-certificates && \
#     addgroup -S node && adduser -S node -G node

# Create application folder and assign rights to the node user
RUN mkdir -p $APP_DIR && chown -R node:node $APP_DIR

# Set the working directory
WORKDIR $APP_DIR

# Copy package.json from the host to image
COPY --chown=node:node package*.json .

# Copy sources and docs from the host to image
COPY --chown=node:node ./src .

# Set the active user
USER node

##
## Stage: Development environment
FROM base AS development

# Define environment variables
# ENV APP_DIR=/opt/app
ENV NODE_ENV=development

# Set the active user
USER node

# Install application modules
RUN npm ci && \
  npm cache clean --force

# Expose port on the host
EXPOSE $NODE_PORT

# Run the application in development mode
CMD ["npm", "run", "start"]

##
## Stage: Production environment
FROM base AS production

# Define environment variables
ENV NODE_ENV=production

# Set the active user
USER node

# Install application modules
RUN npm ci --only=production && \
  npm cache clean --force

# Expose port on the host
EXPOSE $NODE_PORT

# Run the application
# CMD ["node", "src/server.js"]
CMD ["npm", "run", "start"]
