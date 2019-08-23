import { configure } from '@storybook/react';

function loadStories() {
  require('../src/stories/index.stories.js');
  require('../src/stories/statusTimeLine.stories.js');
}

configure(loadStories, module);
