ARG NODE_VERSION=22

FROM node:${NODE_VERSION}-alpine AS build

# Define environment variables
ENV APP_DIR=/opt/build

# Create application folder and assign rights to the node user
RUN mkdir -p $APP_DIR && chown -R node:node $APP_DIR

# Set the working directory
WORKDIR $APP_DIR

# Copy package.json from the host to image
COPY package*.json $APP_DIR/

# Copy sources from the host to image
COPY ./src $APP_DIR/

# Install application modules
RUN npm install && \
  npm cache clean --force

##
## Stage: Development environment
FROM build AS development

# Define environment variables
ENV APP_DIR=/opt/app
ENV NODE_ENV=development
ENV NODE_PORT=8080

# Create application folder and assign rights to the node user
RUN mkdir -p $APP_DIR && chown -R node:node $APP_DIR

# Set the working directory
WORKDIR $APP_DIR

# Set the active user
USER node

# Copy remaining files
COPY --chown=node:node --from=build /opt/build ${APP_DIR}

# Expose port on the host
EXPOSE $NODE_PORT

# Run the application in development mode
CMD ["npm", "run", "dev"]

##
## Stage: Production environment
FROM build AS production

# Define environment variables
ENV APP_DIR=/opt/app
ENV NODE_ENV=production
ENV NODE_PORT=8080

# Create application folder and assign rights to the node user
RUN mkdir -p $APP_DIR && chown -R node:node $APP_DIR

# Set the working directory
WORKDIR $APP_DIR

# Set the active user
USER node

# Copy remaining files
COPY --chown=node:node --from=build /opt/build/node_modules ${APP_DIR}/node_modules
COPY --chown=node:node --from=build /opt/build/src ${APP_DIR}/src

# Expose port on the host
EXPOSE $NODE_PORT

# Run the application
CMD ["npm", "start"]
