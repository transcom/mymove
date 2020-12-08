import React from 'react';
import { text } from '@storybook/addon-knobs';

import SomethingWentWrong from '../shared/SomethingWentWrong';

export default {
  title: 'Components/Something Went Wrong',
};

export const Component = () => (
  <div className="usa-grid">
    <div style={{ textAlign: 'center' }}>
      <SomethingWentWrong
        error={text('SomethingWentWrong.error', 'error')}
        info={text('SomethingWentWrong.info', 'info')}
      />
    </div>
  </div>
);
