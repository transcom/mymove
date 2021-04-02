import React from 'react';

import OrdersInfoForm from './OrdersInfoForm';

const testInitialValues = {
  orders_type: 'PERMANENT_CHANGE_OF_STATION',
  issue_date: '2020-11-08',
  report_by_date: '2020-11-26',
  has_dependents: 'no',
  new_duty_station: {
    address: {
      city: 'Des Moines',
      country: 'US',
      id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
      postal_code: '50309',
      state: 'IA',
      street_address_1: '987 Other Avenue',
      street_address_2: 'P.O. Box 1234',
      street_address_3: 'c/o Another Person',
    },
    address_id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
    affiliation: 'AIR_FORCE',
    created_at: '2020-10-19T17:01:16.114Z',
    id: 'f9299768-16d2-4a13-ae39-7087a58b1f62',
    name: 'Yuma AFB',
    updated_at: '2020-10-19T17:01:16.114Z',
  },
};

export default {
  title: 'Customer Components / Forms / Orders Info Form',
  component: OrdersInfoForm,
  argTypes: {
    onSubmit: { action: 'submit form' },
    onBack: { action: 'go back' },
  },
};

const testProps = {
  initialValues: { orders_type: '', issue_date: '', report_by_date: '', has_dependents: '', new_duty_station: {} },
  ordersTypeOptions: [
    { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
    { key: 'RETIREMENT', value: 'Retirement' },
    { key: 'SEPARATION', value: 'Separation' },
  ],
  currentStation: {},
};

export const EmptyValues = (argTypes) => (
  <OrdersInfoForm {...testProps} onSubmit={argTypes.onSubmit} onBack={argTypes.onBack} />
);

export const PrefillNoDependents = (argTypes) => (
  <OrdersInfoForm
    {...testProps}
    initialValues={testInitialValues}
    onSubmit={argTypes.onSubmit}
    onBack={argTypes.onBack}
  />
);

export const PrefillYesDependents = (argTypes) => (
  <OrdersInfoForm
    {...testProps}
    initialValues={{ ...testInitialValues, has_dependents: 'yes' }}
    onSubmit={argTypes.onSubmit}
    onBack={argTypes.onBack}
  />
);

export const PCSOnly = (argTypes) => (
  <OrdersInfoForm
    {...testProps}
    ordersTypeOptions={[testProps.ordersTypeOptions[0]]}
    onSubmit={argTypes.onSubmit}
    onBack={argTypes.onBack}
  />
);
