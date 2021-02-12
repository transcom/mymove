import React from 'react';
import { action } from '@storybook/addon-actions';

import EulaModal from './index';

export default {
  title: 'Components/Eula Modal',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/d9ad20e6-944c-48a2-bbd2-1c7ed8bc1315?mode=design',
    },
  },
};

const props = {
  acceptTerms: action('clicked'),
  closeModal: action('clicked'),
};

// eslint-disable-next-line react/jsx-props-no-spreading
export const EulaModalStory = () => <EulaModal {...props} />;
