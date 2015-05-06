FROM quay.io/lsqio/iojs

WORKDIR /app
ENTRYPOINT ["/usr/bin/env", "npm", "start", "--"]

ENV PORT 3000
EXPOSE 3000

COPY . /app
RUN ["/usr/bin/env", "npm", "install", "--production", "--loglevel=http"]
