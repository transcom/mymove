import { format } from 'util';

import Enzyme from 'enzyme';
import Adapter from '@wojtekmaj/enzyme-adapter-react-17';
import 'jest-canvas-mock';
import '@testing-library/jest-dom';

import './icons';

// Set up enzyme's React adapter for use in other test files
Enzyme.configure({ adapter: new Adapter() });

// Window/global mocks
global.scrollTo = () => {};

global.performance = {
  now: () => {
    return Date.now();
  },
};

global.console.error = (message, ...args) => {
  if (message.indexOf('Warning: An update to %s inside a test was not wrapped in act(') > -1) {
    const errorText = `
React emitted an error which states that a component re-rendered itself after a test completed.
This indicates that some test that ran earlier within this test suite may not
be correctly or fully capturing the component's behavior.

Using \`async\` functions and \`findBy\` queries can help solve this problem.
(see: https://kentcdodds.com/blog/fix-the-not-wrapped-in-act-warning)
    `;

    throw errorText;
  }

  global.console.log("Jest tests can't emit console.error\n");
  throw format(message, ...args);
};
