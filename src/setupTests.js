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
  global.console.log("Jest tests can't emit console.error\n");
  throw format(message, ...args);
};
