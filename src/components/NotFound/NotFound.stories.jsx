import React from 'react';
import { action } from '@storybook/addon-actions';

import NotFound from './NotFound';

export default {
  title: 'Components / Not Found',
};

export const NotFoundComponent = () => <NotFound handleOnClick={action('clicked')} />;
