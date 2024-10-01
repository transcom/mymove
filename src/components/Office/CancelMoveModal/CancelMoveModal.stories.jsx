import React from 'react';
import { action } from '@storybook/addon-actions';

import CancelMoveModal from './CancelMoveModal';

export default {
  title: 'Office Components/CancelMoveModal',
  component: CancelMoveModal,
};

export const Basic = () => <CancelMoveModal onSubmit={action('Submit')} onClose={action('Cancel')} />;
