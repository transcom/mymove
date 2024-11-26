import Enzyme from 'enzyme';
import Adapter from '@wojtekmaj/enzyme-adapter-react-17';
import 'jest-canvas-mock';
import '@testing-library/jest-dom';
import { jestPreviewConfigure } from 'jest-preview';

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

jestPreviewConfigure({
  autoPreview: false,
});
