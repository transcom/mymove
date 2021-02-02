import React from 'react';

import BypassBlock from './index';

export default {
  title: 'Components/Headers/Bypass Block',
};

export const BypassBlockLink = () => (
  <div>
    <BypassBlock />
    <nav>
      Sample Navigation
      <a href="#">Link</a>
      <a href="#">Link</a>
      <a href="#">Link</a>
    </nav>
  </div>
);
