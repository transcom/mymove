import React from 'react';
import { action } from '@storybook/addon-actions';

import { ValidCACModal } from './ValidCACModal';

export default {
  title: 'Components/Valid CAC Modal',
  component: ValidCACModal,
};

const props = {
  onClose: action('clicked'),
  onSubmit: action('clicked'),
};

export const ValidCAC = () => <ValidCACModal {...props} />;
