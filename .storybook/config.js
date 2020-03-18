import { configure, addDecorator } from '@storybook/react';
import { withInfo } from '@storybook/addon-info';
import 'loki/configure-react';

import './storybook.scss';
import '../src/index.scss';

const req = require.context('../src', true, /\.stories\.jsx?$/);

const loadStories = () => {
  req.keys().forEach(req);
};

addDecorator(withInfo);
configure(loadStories, module);
