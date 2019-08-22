import { configure } from '@storybook/react';

function loadStories() {
  require('../src/stories/index.stories.js');
}

configure(loadStories, module);
