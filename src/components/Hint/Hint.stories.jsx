import React from 'react';

import Hint from './index';

export default {
  title: 'Components/Hint',
  component: Hint,
};

export const Basic = () => (
  <Hint>
    <p>Here is some hint text.</p>
  </Hint>
);

export const MultipleParagraphs = () => (
  <Hint>
    <p>Here is some hint text.</p>
    <p>Here is another paragraph hint text.</p>
  </Hint>
);
