import React from 'react';

import BypassBlock from './index';

export default {
  title: 'Components/Headers/Bypass Block',
};

export const BypassBlockLink = () => (
  <div>
    <BypassBlock />
    <p>Press the Tab key to focus the bypass block.</p>
  </div>
);
