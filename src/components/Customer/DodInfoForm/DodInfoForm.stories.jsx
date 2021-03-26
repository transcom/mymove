import React from 'react';

import DodInfoForm from './DodInfoForm';

export default {
  title: 'Customer Components / Forms / DOD Info Form',
  component: DodInfoForm,
  argTypes: {
    onSubmit: { action: 'submit form' },
    onBack: { action: 'go back' },
  },
};

export const DefaultState = (argTypes) => (
  <DodInfoForm initialValues={{}} onSubmit={argTypes.onSubmit} onBack={argTypes.onBack} />
);

export const WithInitialValues = (argTypes) => (
  <DodInfoForm
    initialValues={{
      affiliation: 'ARMY',
      edipi: '9999999999',
      rank: 'E_2',
    }}
    onSubmit={argTypes.onSubmit}
    onBack={argTypes.onBack}
  />
);
