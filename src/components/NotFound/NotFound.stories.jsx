import React from 'react';
import { action } from '@storybook/addon-actions';

import NotFound from './NotFound';

import { MockRouting } from 'testUtils';

export default {
  title: 'Components / Not Found',
};

export const NotFoundComponent = () => (
  <MockRouting>
    <NotFound handleOnClick={action('clicked')} />
  </MockRouting>
);
