import React from 'react';
import { action } from '@storybook/addon-actions';

import NotFound from './NotFound';

import { MockRouterProvider } from 'testUtils';

export default {
  title: 'Components / Not Found',
};

export const NotFoundComponent = () => (
  <MockRouterProvider>
    <NotFound handleOnClick={action('clicked')} />
  </MockRouterProvider>
);
