# build
FROM node:16.6.0-alpine3.11 AS builder
ARG NODE_ENV=development
ENV NODE_ENV=${NODE_ENV}
WORKDIR /usr/src/app
COPY package*.json /usr/src/app/
RUN npm ci

# run
FROM node:16.6.0-alpine3.11
RUN apk add --no-cache dumb-init
USER node
WORKDIR /usr/src/app
COPY --chown=node:node --from=builder /usr/src/app/node_modules /usr/src/app/node_modules
COPY --chown=node:node . /usr/src/app
EXPOSE 5000
CMD [ "dumb-init", "npm", "start" ]

