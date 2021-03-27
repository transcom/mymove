import React from 'react';

import CurrentDutyStationForm from './CurrentDutyStationForm';

export default {
  title: 'Customer Components / Forms / Current Duty Station Form',
  component: CurrentDutyStationForm,
  argTypes: {
    onSubmit: { action: 'submit form' },
    onBack: { action: 'go back' },
    initialValues: {
      current_station: {},
    },
  },
};

export const DefaultState = (argTypes) => (
  <CurrentDutyStationForm
    initialValues={argTypes.initialValues}
    onSubmit={argTypes.onSubmit}
    onBack={argTypes.onBack}
  />
);
