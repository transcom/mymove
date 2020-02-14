import { configure } from '@storybook/react';

import 'uswds';
import 'uswds/dist/css/uswds.css';

function loadStories() {
  require('../src/stories/index.stories.jsx');
  require('../src/stories/statusTimeLine.stories.jsx');
}

configure(loadStories, module);
