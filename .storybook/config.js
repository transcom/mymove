import { configure } from '@storybook/react';

// automatically import all files ending in *.stories.js
//const req = require.context('../src/stories', true, /\.stories\.js$/);
function loadStories() {
  require('../src/stories/index.stories.js');
  // This does not work in IE 11
  //req.keys().forEach(filename => req(filename));
}

configure(loadStories, module);
