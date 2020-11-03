import React from 'react';
import { Formik } from 'formik';

import OrdersInfoForm from './OrdersInfoForm';

export default {
  title: 'Customer Components | Forms / OrdersInfoForm',
  component: OrdersInfoForm,
  decorators: [(storyFn) => <Formik>{storyFn()}</Formik>],
};

export const EmptyValues = () => <OrdersInfoForm />;
