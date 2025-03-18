import React from 'react';
import { Formik } from 'formik';

import MaskedTextField from './MaskedTextField';

import { Form } from 'components/form/Form';

export default {
  title: 'Components/MaskedTextFields',
};

const labelHint =
  'This TAC does not appear in TGET, so it might not be valid. Make sure it matches what&apos;s on the orders before you continue.';

export const MaskedTextFieldDefaultState = () => (
  <Formik initialValues={{}}>
    {() => (
      <Form>
        <MaskedTextField
          id="input-type-text"
          label="Text input label"
          hint={labelHint}
          name="input-type-text"
          type="text"
        />
      </Form>
    )}
  </Formik>
);
export const MaskedTextFieldDisabledState = () => (
  <Formik initialValues={{}}>
    {() => (
      <Form>
        <MaskedTextField
          id="input-type-text"
          label="Text input label"
          hint={labelHint}
          name="input-type-text"
          type="text"
          isDisabled
        />
      </Form>
    )}
  </Formik>
);
export const MaskedTextFieldWithWarning = () => (
  <Formik initialValues={{}}>
    {() => (
      <Form>
        <MaskedTextField
          id="input-type-text"
          label="Text input label"
          name="input-type-text"
          type="text"
          warning="This TAC does not appear in TGET, so it might not be valid. Make sure it matches what's on the orders before you continue."
        />
      </Form>
    )}
  </Formik>
);

export const MaskedTextFieldWithError = () => (
  <Formik initialValues={{}}>
    {() => (
      <Form>
        <MaskedTextField
          id="input-type-text"
          label="Text input label"
          name="input-type-text"
          type="text"
          validationStatus="error"
          errorMessage="Helpful error message"
          error
        />
      </Form>
    )}
  </Formik>
);

export const MaskedTextFieldWithOptionalTag = () => (
  <Formik initialValues={{}}>
    {() => (
      <Form>
        <MaskedTextField id="input-type-text" label="Text input label" name="input-type-text" type="text" optional />
      </Form>
    )}
  </Formik>
);
