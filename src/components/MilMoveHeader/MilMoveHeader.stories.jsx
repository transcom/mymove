import React from 'react';
import { action } from '@storybook/addon-actions';

import MilMoveHeader from './index';

export default {
  title: 'Components/Headers/MilMove Header',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/d9ad20e6-944c-48a2-bbd2-1c7ed8bc1315?mode=design',
    },
  },
};

const props = {
  officeUser: { last_name: 'Baker', first_name: 'Riley' },
  handleLogout: action('clicked'),
};

export const Milmove = () => (
  // eslint-disable-next-line react/jsx-props-no-spreading
  <MilMoveHeader {...props}>
    <a href="#">Navigation Link</a>
    <a href="#">Navigation Link</a>
    <a href="#">Navigation Link</a>
  </MilMoveHeader>
);
