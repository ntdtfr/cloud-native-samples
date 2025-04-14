# syntax=docker/dockerfile:1

ARG NODE_VERSION=22

##
## Stage 1: Development dependencies and build
FROM node:${NODE_VERSION}-alpine AS builder

# Install git for Husky
RUN apk add --no-cache git

# Set working directory
WORKDIR /opt/app

# Add package files and install dependencies
COPY package.json package-lock.json ./

# Install npm version 11.3.0
RUN npm install -g npm@11.3.0 && \
  npm install -g babel-cli

# Install all dependencies including dev dependencies
# Use ci instead of install for more reliable builds
RUN npm ci

# Set up Git config to avoid Husky errors
RUN git config --global init.defaultBranch main && \
    git config --global user.email "ntran@ntdt.fr" && \
    git config --global user.name "nntran"

# Initialize git repository to allow Husky installation
RUN git init

# Copy source code
COPY . .

# # Run prepare script which installs Husky
# RUN npm run prepare

# Build application (adjust according to your build script)
RUN npm run build

##
## Stage 2: Production dependencies only
FROM node:${NODE_VERSION}-alpine AS dependencies

WORKDIR /opt/app

COPY package.json package-lock.json ./

# Install npm version 11.3.0
RUN npm install -g npm@11.3.0

# Install only production dependencies
RUN npm ci --only=production

##
## Stage 3: Final image
FROM node:${NODE_VERSION}-alpine AS production

# Set non-root user for security
USER node

# Set environment variables
ENV NODE_ENV=production
ENV PORT=3000

# Create app directory and set correct ownership
WORKDIR /opt/app

# Copy build artifacts and dependencies from previous stages
COPY --from=builder --chown=node:node /opt/app/dist ./dist
COPY --from=dependencies --chown=node:node /opt/app/node_modules ./node_modules
COPY --chown=node:node package.json ./

# Expose port
EXPOSE ${PORT}

# Use healthcheck to verify container health
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:${PORT}/health || exit 1

# Use specific command to run the application
CMD ["node", "dist/server.js"]
