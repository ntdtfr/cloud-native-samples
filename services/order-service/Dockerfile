ARG NODE_VERSION=22

FROM node:${NODE_VERSION}-alpine AS base

# Define environment variables
ENV APP_DIR=/opt/app

# Set the working directory
WORKDIR $APP_DIR

# Copy package.json from the host to image
COPY package*.json $APP_DIR/

# Copy sources from the host to image
COPY . $APP_DIR/

##
## Stage: Development environment
FROM base AS development

# Define environment variables
ENV APP_DIR=/opt/app
ENV NODE_ENV=development
ENV NODE_PORT=3000

# Create application folder and assign rights to the node user
RUN mkdir -p $APP_DIR && chown -R node:node $APP_DIR

# Set the working directory
WORKDIR $APP_DIR

# Install application modules
RUN npm install && \
  npm cache clean --force

# Set the active user
USER node

# Copy remaining files
COPY --chown=node:node --from=base /opt/app ${APP_DIR}

# Expose port on the host
EXPOSE $NODE_PORT

# Run the application in development mode
CMD ["npm", "run", "start"]

##
## Stage: Production environment
FROM base AS production

# Define environment variables
ENV APP_DIR=/opt/app
ENV NODE_ENV=production
ENV NODE_PORT=3000

# Create application folder and assign rights to the node user
RUN mkdir -p $APP_DIR && chown -R node:node $APP_DIR

# Set the working directory
WORKDIR $APP_DIR

# Install application modules
RUN npm ci --only=production && \
  npm cache clean --force

# Set the active user
USER node

# Copy remaining files
COPY --chown=node:node --from=build /opt/build/node_modules ${APP_DIR}/node_modules
COPY --chown=node:node --from=build /opt/build/src ${APP_DIR}/src

# Expose port on the host
EXPOSE $NODE_PORT

# Run the application
CMD ["node", "src/server.js"]
