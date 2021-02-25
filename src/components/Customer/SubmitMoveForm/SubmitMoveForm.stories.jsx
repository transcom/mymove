import React from 'react';

import SubmitMoveForm from './SubmitMoveForm';

export default {
  title: 'Customer Components / Forms / SubmitMoveForm',
  component: SubmitMoveForm,
  argTypes: {
    onSubmit: { action: 'submit form' },
    onPrint: { action: 'print page' },
  },
};

export const DefaultState = (argTypes) => <SubmitMoveForm onSubmit={argTypes.onSubmit} onPrint={argTypes.onPrint} />;

export const WithServerError = (argTypes) => (
  <SubmitMoveForm onSubmit={argTypes.onSubmit} onPrint={argTypes.onPrint} error />
);
