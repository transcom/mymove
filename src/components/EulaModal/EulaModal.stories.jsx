import React from 'react';
import { action } from '@storybook/addon-actions';

import { EulaModal } from './index';

export default {
  title: 'Components/Eula Modal',
  component: EulaModal,
};

const props = {
  acceptTerms: action('clicked'),
  closeModal: action('clicked'),
};

// eslint-disable-next-line react/jsx-props-no-spreading
export const EulaModalStory = () => <EulaModal {...props} />;
