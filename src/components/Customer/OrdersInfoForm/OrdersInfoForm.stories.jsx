/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { Formik } from 'formik';

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
  title: 'Customer Components / Forms / OrdersInfoForm',
  component: OrdersInfoForm,
  // TODO - Story-level decorators not working, maybe after we update Storybook?
  /* decorators: [
    (Story) => (
      <Formik initialValues={{ ...testInitialValues, has_dependents: 'no' }}>
        <Story />
      </Formik>
    ),
  ], */
};

const testProps = {
  ordersTypeOptions: [
    { key: 'PERMANENT_CHANGE_OF_STATION', value: 'Permanent Change Of Station (PCS)' },
    { key: 'RETIREMENT', value: 'Retirement' },
    { key: 'SEPARATION', value: 'Separation' },
  ],
};

export const EmptyValues = () => (
  <Formik>
    <OrdersInfoForm {...testProps} />
  </Formik>
);

export const PrefillNoDependents = () => (
  <Formik initialValues={testInitialValues}>
    <OrdersInfoForm {...testProps} />
  </Formik>
);

export const PrefillYesDependents = () => (
  <Formik initialValues={{ ...testInitialValues, has_dependents: 'yes' }}>
    <OrdersInfoForm {...testProps} />
  </Formik>
);

export const PCSOnly = () => (
  <Formik>
    <OrdersInfoForm ordersTypeOptions={[testProps.ordersTypeOptions[0]]} />
  </Formik>
);
