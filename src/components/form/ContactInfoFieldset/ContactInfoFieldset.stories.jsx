import React from 'react';
import { action } from '@storybook/addon-actions';
import { Formik } from 'formik';

import { ContactInfoFieldset } from './index';

export default {
  title: 'Components/ContactInfoFieldset',
};

const props = {
  name: 'contactInfoFieldset',
  onChangePreferEmail: action('clicked'),
  onChangePreferPhone: action('clicked'),
};

export const ContactInfoFieldsetStory = () => (
  <Formik>
    {/* eslint-disable-next-line react/jsx-props-no-spreading */}
    <ContactInfoFieldset {...props} />
  </Formik>
);
