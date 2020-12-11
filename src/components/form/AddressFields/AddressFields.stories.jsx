import React from 'react';
import { Formik } from 'formik';

import { AddressFields } from './AddressFields';

export default {
  title: 'Components/AddressFields',
};

export const Basic = () => (
  <Formik>
    <AddressFields />
  </Formik>
);
