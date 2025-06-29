import React from 'react';

import OrdersInfoForm from './OrdersInfoForm';

import { ORDERS_TYPE } from 'constants/orders';
import { MockProviders } from 'testUtils';

const testInitialValues = {
  orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION,
  issue_date: '2020-11-08',
  report_by_date: '2020-11-26',
  has_dependents: 'no',
  new_duty_location: {
    address: {
      city: 'Des Moines',
      country: 'US',
      id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
      postalCode: '50309',
      state: 'IA',
      streetAddress1: '987 Other Avenue',
      streetAddress2: 'P.O. Box 1234',
      streetAddress3: 'c/o Another Person',
    },
    address_id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
    affiliation: 'AIR_FORCE',
    created_at: '2020-10-19T17:01:16.114Z',
    id: 'f9299768-16d2-4a13-ae39-7087a58b1f62',
    name: 'Yuma AFB',
    updated_at: '2020-10-19T17:01:16.114Z',
  },
  grade: 'E-1',
  origin_duty_location: {
    address: {
      city: 'Des Moines',
      country: 'US',
      id: 'a4b30b99-4e82-48a6-b736-01662b499d6a',
      postalCode: '50309',
      state: 'IA',
      streetAddress1: '987 Other Avenue',
      streetAddress2: 'P.O. Box 1234',
      streetAddress3: 'c/o Another Person',
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
  initialValues: {
    orders_type: '',
    issue_date: '',
    report_by_date: '',
    has_dependents: '',
    new_duty_location: {},
    grade: '',
    origin_duty_location: {},
  },
  ordersTypeOptions: [
    { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
    { key: 'LOCAL_MOVE', value: 'Local Move' },
    { key: 'RETIREMENT', value: 'Retirement' },
    { key: 'SEPARATION', value: 'Separation' },
  ],
};

export const EmptyValues = (argTypes) => (
  <MockProviders>
    <OrdersInfoForm {...testProps} onSubmit={argTypes.onSubmit} onBack={argTypes.onBack} />
  </MockProviders>
);

export const PrefillNoDependents = (argTypes) => (
  <MockProviders>
    <OrdersInfoForm
      {...testProps}
      initialValues={testInitialValues}
      onSubmit={argTypes.onSubmit}
      onBack={argTypes.onBack}
    />
  </MockProviders>
);

export const PrefillYesDependents = (argTypes) => (
  <MockProviders>
    <OrdersInfoForm
      {...testProps}
      initialValues={{ ...testInitialValues, has_dependents: 'yes' }}
      onSubmit={argTypes.onSubmit}
      onBack={argTypes.onBack}
    />
  </MockProviders>
);

export const PCSOnly = (argTypes) => (
  <MockProviders>
    <OrdersInfoForm
      {...testProps}
      ordersTypeOptions={[testProps.ordersTypeOptions[0]]}
      onSubmit={argTypes.onSubmit}
      onBack={argTypes.onBack}
    />
  </MockProviders>
);
