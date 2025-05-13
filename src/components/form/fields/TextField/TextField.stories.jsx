import React from 'react';
import { Formik } from 'formik';

import TextField from './TextField';

import { Form } from 'components/form/Form';

export default {
  title: 'Components/TextFields',
};

const labelHint =
  'This TAC does not appear in TGET, so it might not be valid. Make sure it matches what&apos;s on the orders before you continue.';

export const TextFieldDefaultState = () => (
  <Formik initialValues={{}}>
    {() => (
      <Form>
        <TextField id="input-type-text" label="Text input label" hint={labelHint} name="input-type-text" type="text" />
      </Form>
    )}
  </Formik>
);
export const TextFieldDisabledState = () => (
  <Formik initialValues={{}}>
    {() => (
      <Form>
        <TextField
          id="input-type-text"
          label="Text input label"
          hint={labelHint}
          name="input-type-text"
          isDisabled
          type="text"
        />
      </Form>
    )}
  </Formik>
);
export const TextFieldWithWarning = () => (
  <Formik initialValues={{}}>
    {() => (
      <Form>
        <TextField
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

export const TextFieldWithError = () => (
  <Formik initialValues={{}}>
    {() => (
      <Form>
        <TextField
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

export const TextFieldWithOptionalTag = () => (
  <Formik initialValues={{}}>
    {() => (
      <Form>
        <TextField id="input-type-text" label="Text input label" name="input-type-text" type="text" optional />
      </Form>
    )}
  </Formik>
);

export const TextFieldMultiLine = () => (
  <Formik initialValues={{}}>
    {() => (
      <Form>
        <TextField
          id="input-type-text"
          label="Text input label"
          name="input-type-text"
          type="text"
          display="textarea"
        />
      </Form>
    )}
  </Formik>
);

export const TextFieldReadOnly = () => (
  <Formik
    initialValues={{
      'input-type-text': 'Text Value',
    }}
  >
    {() => (
      <Form>
        <TextField
          id="input-type-text"
          label="Text input label"
          name="input-type-text"
          type="text"
          display="readonly"
        />
      </Form>
    )}
  </Formik>
);

export const TextFieldWithRequiredAsterisk = () => (
  <Formik initialValues={{}}>
    {() => (
      <Form>
        <TextField
          id="input-type-text"
          label="Text input label"
          hint={labelHint}
          name="input-type-text"
          type="text"
          showRequiredAsterisk
        />
      </Form>
    )}
  </Formik>
);

export const TextFieldWithPrefix = () => (
  <Formik initialValues={{}}>
    {() => (
      <Form>
        <TextField
          label="Prefix example"
          id="prefixExample"
          name="prefixExample"
          prefix="PREFIX EXAMPLE:"
          showRequiredAsterisk
        />
      </Form>
    )}
  </Formik>
);
