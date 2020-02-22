import { configure, addDecorator } from '@storybook/react';
import { withInfo } from '@storybook/addon-info';
import 'loki/configure-react';

import 'uswds';
import 'uswds/dist/css/uswds.css';

function loadStories() {
  require('../src/stories/index.stories.jsx');
  require('../src/stories/statusTimeLine.stories.jsx');
  require('../src/stories/tabNav.stories.jsx');
}

addDecorator(withInfo);
configure(loadStories, module);
