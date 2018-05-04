import Enzyme from 'enzyme';
import Adapter from 'enzyme-adapter-react-16';

// Set up enzyme's React adapter for use in other test files
Enzyme.configure({ adapter: new Adapter() });
global.performance = {
  now: function() {
    return Date.now();
  },
};
