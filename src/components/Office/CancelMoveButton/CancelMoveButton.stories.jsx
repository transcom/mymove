import React from 'react';
import { action } from '@storybook/addon-actions';

import CancelMoveButton from './CancelMoveButton';

export default {
  title: 'Office Components/CancelMoveButton',
  component: CancelMoveButton,
};

export const MoveNotCanceled = () => (
  <div className="officeApp">
    <CancelMoveButton onClick={action('Click')} />
  </div>
);
