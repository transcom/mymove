import { configure } from '@storybook/react';

import 'uswds';
import 'uswds/dist/css/uswds.css';

function loadStories() {
  require('../src/stories/index.stories.js');
  require('../src/stories/statusTimeLine.stories.js');
  require('../src/stories/dateAndLocation.stories.js');
}

configure(loadStories, module);
