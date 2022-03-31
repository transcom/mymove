import React from 'react';

import CurrentDutyLocationForm from './CurrentDutyLocationForm';

export default {
  title: 'Customer Components / Forms / Current Duty Location Form',
  component: CurrentDutyLocationForm,
  argTypes: {
    onSubmit: { action: 'submit form' },
    onBack: { action: 'go back' },
  },
};

export const DefaultState = (argTypes) => (
  <CurrentDutyLocationForm initialValues={{}} onSubmit={argTypes.onSubmit} onBack={argTypes.onBack} />
);

export const InitialValues = (argTypes) => (
  <CurrentDutyLocationForm
    initialValues={{
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
    onBack={argTypes.onBack}
  />
);

export const Error = (argTypes) => (
  <CurrentDutyLocationForm
    initialValues={{
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
    newDutyLocation={{
      address: {
        city: 'Los Angeles',
        state: 'CA',
        postalCode: '90245',
      },
      id: 'testId',
    }}
    onSubmit={argTypes.onSubmit}
    onBack={argTypes.onBack}
  />
);
