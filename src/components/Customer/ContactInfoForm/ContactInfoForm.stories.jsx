import React from 'react';

import ContactInfoForm from './index';

export default {
  title: 'Customer Components / Forms/ Contact Info Form',
};

const initialValues = {};
export const DefaultState = () => <ContactInfoForm initialValues={initialValues} />;
