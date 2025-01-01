FROM node:23.1.0-alpine

# define environment variables
ENV APP_DIR=/opt/video-streaming
ENV NODE_ENV=production
ENV NODE_PORT=3000

# Create application folder and assign rights to the node user
RUN mkdir -p $APP_DIR && chown -R node:node $APP_DIR

# Set the working directory
WORKDIR $APP_DIR

# Set the active user
USER node

# Copy package.json from the host
COPY --chown=node:node package*.json $APP_DIR/

# Install application modules
RUN npm install --only=production && \
  npm cache clean --force

# Copy remaining files
#COPY --chown=node:node . .
COPY --chown=node:node ./src ./src
COPY --chown=node:node ./videos ./videos

# Expose port on the host
EXPOSE $NODE_PORT

# Run the application
CMD ["npm", "start"]
