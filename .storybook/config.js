import { configure, addDecorator } from '@storybook/react';
import { withInfo } from '@storybook/addon-info';
import 'loki/configure-react';

import './storybook.scss';
import '../src/index.scss';

function loadStories() {
  require('../src/stories/index.stories.jsx');
  require('../src/stories/statusTimeLine.stories.jsx');
  require('../src/stories/tabNav.stories.jsx');
}

addDecorator(withInfo);
configure(loadStories, module);
