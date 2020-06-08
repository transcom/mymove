import { configure, addDecorator } from '@storybook/react';
import { withInfo } from '@storybook/addon-info';

import './storybook.scss';
import '../src/index.scss';
import '../src/ghc_index.scss';
import { detectIE11 } from '../src/shared/utils';

if (!detectIE11()) {
  // eslint-disable-next-line no-unused-expressions
  import('loki/configure-react');
}

const req = require.context('../src', true, /\.stories\.jsx?$/);

const loadStories = () => {
  req.keys().forEach(req);
};

addDecorator(withInfo);
configure(loadStories, module);
