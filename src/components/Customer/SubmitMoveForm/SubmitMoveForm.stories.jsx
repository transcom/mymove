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
    onSubmit={argTypes.onSubmit}
    certificationText={completeCertificationText}
    onPrint={argTypes.onPrint}
  />
);

export const WithServerError = (argTypes) => (
  <SubmitMoveForm
    onSubmit={argTypes.onSubmit}
    onPrint={argTypes.onPrint}
    certificationText={completeCertificationText}
    error
  />
);

export const LoadingCertificationText = (argTypes) => (
  <SubmitMoveForm onSubmit={argTypes.onSubmit} onPrint={argTypes.onPrint} />
);
