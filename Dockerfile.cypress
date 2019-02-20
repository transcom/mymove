FROM cypress/base:10


# We use `npm install` here because you can't selectively `yarn install` packages and we don't care about a lockfile and
# `yarn add` won't read from our existing package.json
COPY package.json /package.json
RUN npm install --save-dev mocha mocha-multi-reporters mocha-junit-reporter
RUN npm install --save-dev --silent cypress

COPY cypress/ /cypress
COPY cypress.json /cypress.json
COPY mocha-reporter-config.json /mocha-reporter-config.json

ENTRYPOINT ["/node_modules/.bin/cypress"]
CMD ["run"]
