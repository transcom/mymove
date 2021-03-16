import React from 'react';

import SubmitMoveForm from './SubmitMoveForm';

import { completeCertificationText } from 'scenes/Legalese/legaleseText';

export default {
  title: 'Customer Components / Forms / SubmitMoveForm',
  component: SubmitMoveForm,
  argTypes: {
    onSubmit: { action: 'submit form' },
    onPrint: { action: 'print page' },
  },
};

export const DefaultState = (argTypes) => (
  <SubmitMoveForm
    initialValues={{ signature: '', date: '2021-01-20' }}
    onSubmit={argTypes.onSubmit}
    certificationText={completeCertificationText}
    onPrint={argTypes.onPrint}
  />
);

export const WithServerError = (argTypes) => (
  <SubmitMoveForm
    initialValues={{ signature: '', date: '2021-01-20' }}
    onSubmit={argTypes.onSubmit}
    onPrint={argTypes.onPrint}
    certificationText={completeCertificationText}
    error
  />
);

export const LoadingCertificationText = (argTypes) => (
  <SubmitMoveForm
    initialValues={{ signature: '', date: '2021-01-20' }}
    onSubmit={argTypes.onSubmit}
    onPrint={argTypes.onPrint}
  />
);
