import React from 'react';

import Callout from './index';

export default {
  title: 'Components/Callout',
};

export const Component = () => (
  <Callout>
    Examples
    <ul>
      <li>Things that might need special handling</li>
      <li>Access info for a location</li>
      <li>Weapons or alcohol</li>
    </ul>
  </Callout>
);
