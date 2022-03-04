import React from 'react';

import ServiceInfoForm from './ServiceInfoForm';

export default {
  title: 'Customer Components / Forms / Service Info Form',
  component: ServiceInfoForm,
  argTypes: {
    onSubmit: { action: 'submit form' },
    onCancel: { action: 'cancel' },
  },
};

export const DefaultState = (argTypes) => (
  <ServiceInfoForm
    initialValues={{
      first_name: '',
      middle_name: '',
      last_name: '',
      suffix: '',
      affiliation: '',
      edipi: '',
      rank: '',
      current_location: {},
    }}
    onSubmit={argTypes.onSubmit}
    onCancel={argTypes.onCancel}
  />
);

export const WithInitialValues = (argTypes) => (
  <ServiceInfoForm
    initialValues={{
      first_name: 'Leo',
      middle_name: 'Star',
      last_name: 'Spaceman',
      suffix: 'Mr.',
      affiliation: 'ARMY',
      edipi: '9999999999',
      rank: 'E_2',
      current_location: {
        address: {
          city: 'Los Angeles',
          state: 'CA',
          postalCode: '90245',
        },
        name: 'Los Angeles AFB',
        id: 'testId',
      },
    }}
    onSubmit={argTypes.onSubmit}
    onCancel={argTypes.onCancel}
  />
);
