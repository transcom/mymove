/* eslint-disable import/no-extraneous-dependencies, no-console */
const commander = require('commander');
const Conf = require('conf');
const debug = require('debug')('debug');

const { schema } = require('./constants');
const { totalDuration } = require('./scenarios');

const config = new Conf({
  cwd: __dirname, // saves the config file to the same dir of this script
  schema,
});

const program = new commander.Command();

const runAction = async ({ scenario, measurementType, host, verbose, saveReports }) => {
  if (verbose) {
    debug.enabled = true;
  }
  console.log(`Running scenario ${scenario} with measurement ${measurementType}`);

  if (measurementType === 'total-duration') {
    const elapsedTimeResults = await totalDuration(host, config.store, debug, saveReports).catch(() => {
      process.exit(1);
    });

    console.table(elapsedTimeResults);
  } else if (measurementType === 'network-comparison') {
    const results = {};
    // await cannot be used inside of a forEach loop
    // eslint-disable-next-line no-restricted-syntax
    for (const speed of ['fast', 'medium', 'slow']) {
      const configStore = {};
      Object.assign(configStore, config.store, { network: speed });
      console.log(`Running network test with ${speed} profile`);

      // Running these tests in parallel would likely skew the results
      // eslint-disable-next-line no-await-in-loop
      const elapsedTimeResults = await totalDuration(host, configStore, debug, saveReports).catch(() => {
        process.exit(1);
      });
      results[`${speed}`] = elapsedTimeResults;
    }

    console.table(results);
  }
};

program
  .command('run', { isDefault: true })
  .description('runs a benchmark test')
  .addOption(
    new commander.Option('-s, --scenario <scenario>', 'scenario is the page or workflow being tested')
      .default('too-orders-document-viewer')
      .choices(['too-orders-document-viewer'])
      .makeOptionMandatory(),
  )
  .addOption(
    new commander.Option('-m --measurement-type <type>', 'specifies the kind of performance output metrics to measure')
      .default('total-duration')
      .choices(['total-duration', 'network-comparison']),
  )
  .addOption(
    new commander.Option('-h --host <host>', 'base host url to use including port').default('http://officelocal:3000'),
  )
  .addOption(new commander.Option('-v --verbose', 'shows verbose debugging info').default(false))
  .addOption(
    new commander.Option(
      '-r --save-reports',
      'save the reports from lighthouse and performance trace json files',
    ).default(false),
  )
  .action(runAction);

program.parse(process.argv);
