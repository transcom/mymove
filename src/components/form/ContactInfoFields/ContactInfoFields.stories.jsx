import React from 'react';
import { Formik } from 'formik';

import { ContactInfoFields } from './ContactInfoFields';

export default {
  title: 'Components|ContactInfoFields',
};

export const Basic = () => (
  <Formik>
    <ContactInfoFields />
  </Formik>
);
