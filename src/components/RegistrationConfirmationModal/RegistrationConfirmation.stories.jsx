import React from 'react';
import { action } from '@storybook/addon-actions';

import { RegistrationConfirmationModal } from './RegistrationConfirmationModal';

export default {
  title: 'Components/Registration Confirmation Modal',
  component: RegistrationConfirmationModal,
};

const props = {
  onSubmit: action('clicked'),
};

export const RegistrationConfirmation = () => <RegistrationConfirmationModal {...props} />;
