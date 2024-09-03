import React from 'react';
import { Provider } from 'react-redux';

import CustomerContactInfoForm from './CustomerContactInfoForm';

import { configureStore } from 'shared/store';

export default {
  title: 'Office Components / Forms/ Customer Contact Info Form',
  component: CustomerContactInfoForm,
  argTypes: {
    onSubmit: { action: 'submitted on Next' },
    onBack: { action: 'to previous page on Back' },
  },
};

const mockStore = configureStore({});

export const DefaultState = (argTypes) => (
  <Provider store={mockStore.store}>
    <div className="officeApp">
      <CustomerContactInfoForm initialValues={{}} onSubmit={argTypes.onSubmit} onBack={argTypes.onBack} />
    </div>
  </Provider>
);

export const WithInitialValues = (argTypes) => (
  <Provider store={mockStore.store}>
    <div className="officeApp">
      <CustomerContactInfoForm
        initialValues={{
          firstName: 'Leo',
          middleName: 'Star',
          lastName: 'Spaceman',
          suffix: 'Mr.',
          customerTelephone: '555-555-5555',
          customerEmail: 'test@sample.com',
          customerAddress: {
            streetAddress1: '235 Prospect Valley Road SE',
            streetAddress2: 'Apt. 3B',
            city: 'El Paso',
            state: 'TX',
            postalCode: '79912',
          },
          name: 'Leo Spaceman',
          telephone: '555-555-5555',
          email: 'test@sample.com',
        }}
        onNextClick={argTypes.onSubmit}
        onCancelClick={argTypes.onBack}
      />
    </div>
  </Provider>
);
