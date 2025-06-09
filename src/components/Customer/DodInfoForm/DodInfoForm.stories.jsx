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
  <DodInfoForm initialValues={{ edipi: '9999999999' }} onSubmit={argTypes.onSubmit} onBack={argTypes.onBack} />
);

export const WithInitialValues = (argTypes) => (
  <DodInfoForm
    initialValues={{
      affiliation: 'ARMY',
      edipi: '9999999999',
      grade: 'E-2',
    }}
    onSubmit={argTypes.onSubmit}
    onBack={argTypes.onBack}
  />
);

export const CoastGuardCustomer = (argTypes) => (
  <DodInfoForm
    initialValues={{
      affiliation: 'COAST_GUARD',
      edipi: '9999999999',
      grade: 'E-2',
      emplid: '1263456',
    }}
    onSubmit={argTypes.onSubmit}
    onBack={argTypes.onBack}
  />
);
