import React from 'react';
import { Provider } from 'react-redux';

import ContactInfoForm from './ContactInfoForm';

import { configureStore } from 'shared/store';

export default {
  title: 'Office Components/ContactInfoForm',
  component: ContactInfoForm,
  argTypes: {
    onSubmit: {
      action: 'save form',
    },
    onCancel: {
      action: 'cancel',
    },
  },
};

const fakeData = {
  firstName: 'John',
  middleName: 'F',
  lastName: 'Doe',
  telephone: '915-555-2942',
  email: 'mm@example.com',
};

const mockStore = configureStore({});

export const DefaultState = (argTypes) => (
  <Provider store={mockStore.store}>
    <ContactInfoForm
      initialValues={{
        firstName: '',
        middleName: '',
        lastName: '',
        email: '',
        telephone: '',
      }}
      onCancel={argTypes.onCancel}
      onSubmit={argTypes.onSubmit}
    />
  </Provider>
);

export const WithInitialValues = (argTypes) => (
  <Provider store={mockStore.store}>
    <ContactInfoForm initialValues={fakeData} onCancel={argTypes.onCancel} onSubmit={argTypes.onSubmit} />
  </Provider>
);
