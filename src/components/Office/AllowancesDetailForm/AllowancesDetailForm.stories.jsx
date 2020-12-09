import React from 'react';
import { withKnobs } from '@storybook/addon-knobs';
import * as Yup from 'yup';
import { Formik } from 'formik';

import AllowancesDetailForm from 'components/Office/AllowancesDetailForm/AllowancesDetailForm';

export default {
  title: 'TOO/TIO Components|AllowancesDetailForm',
  component: AllowancesDetailForm,
  decorators: [
    withKnobs,
    (Story) => (
      <div style={{ padding: `20px`, background: `#f0f0f0` }}>
        <Story />
      </div>
    ),
  ],
};

export const EmptyValues = () => (
  <Formik
    initialValues={{
      authorizedWeight: '0',
    }}
  >
    <form>
      <AllowancesDetailForm />
    </form>
  </Formik>
);

export const InitialValues = () => {
  return (
    <>
      <Formik
        initialValues={{
          authorizedWeight: '8000',
        }}
        validationSchema={Yup.object({
          authorizedWeight: Yup.number().required('Required'),
        })}
      >
        <form>
          <AllowancesDetailForm />
        </form>
      </Formik>
    </>
  );
};
