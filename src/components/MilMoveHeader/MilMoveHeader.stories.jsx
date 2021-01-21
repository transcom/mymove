/* eslint-disable react/jsx-props-no-spreading */

import React from 'react';
import { action } from '@storybook/addon-actions';
import { MemoryRouter } from 'react-router';

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
  customer: { last_name: 'Baker', first_name: 'Riley', dodID: '999999999' },
  handleLogout: action('clicked'),
};

export const Milmove = () => (
  <MemoryRouter>
    <MilMoveHeader {...props}>
      {' '}
      <span>
        <a href="#">Navigation Link</a>
      </span>
      <span>
        <a href="#">Navigation Link</a>
      </span>
      <span>
        <a href="#">Navigation Link</a>
      </span>
    </MilMoveHeader>
  </MemoryRouter>
);
