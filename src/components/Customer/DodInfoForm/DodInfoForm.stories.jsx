/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { Formik } from 'formik';

import DodInfoForm from './DodInfoForm';

export default {
  title: 'Customer Components | Forms / DOD Info Form',
  component: DodInfoForm,
};

export const EmptyValues = () => (
  <Formik>
    <DodInfoForm />
  </Formik>
);
