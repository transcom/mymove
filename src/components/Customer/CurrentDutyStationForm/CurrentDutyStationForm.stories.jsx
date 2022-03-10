import React from 'react';

import CurrentDutyStationForm from './CurrentDutyStationForm';

export default {
  title: 'Customer Components / Forms / Current Duty Location Form',
  component: CurrentDutyStationForm,
  argTypes: {
    onSubmit: { action: 'submit form' },
    onBack: { action: 'go back' },
  },
};

export const DefaultState = (argTypes) => (
  <CurrentDutyStationForm initialValues={{}} onSubmit={argTypes.onSubmit} onBack={argTypes.onBack} />
);

export const InitialValues = (argTypes) => (
  <CurrentDutyStationForm
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
  <CurrentDutyStationForm
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
