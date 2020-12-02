import Enzyme from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';
import 'jest-canvas-mock';

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
