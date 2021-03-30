import React from 'react';

import CurrentDutyStationForm from './CurrentDutyStationForm';

export default {
  title: 'Customer Components / Forms / Current Duty Station Form',
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
      current_station: {
        address: {
          city: 'Los Angeles',
          state: 'CA',
          postal_code: '90245',
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
      current_station: {
        address: {
          city: 'Los Angeles',
          state: 'CA',
          postal_code: '90245',
        },
        name: 'Los Angeles AFB',
        id: 'testId',
      },
    }}
    newDutyStation={{
      address: {
        city: 'Los Angeles',
        state: 'CA',
        postal_code: '90245',
      },
      id: 'testId',
    }}
    onSubmit={argTypes.onSubmit}
    onBack={argTypes.onBack}
  />
);
