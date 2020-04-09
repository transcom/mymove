import React from 'react';

import { useFormikContext } from 'formik';
import { Form as UswdsForm } from '@trussworks/react-uswds';

export const Form = (props) => {
  const { handleReset, handleSubmit } = useFormikContext();
  // eslint-disable-next-line react/jsx-props-no-spreading
  return <UswdsForm onSubmit={handleSubmit} onReset={handleReset} {...props} />;
};

export default Form;
