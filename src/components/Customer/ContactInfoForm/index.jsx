import React from 'react';
import { Formik } from 'formik';

import { CustomerContactInfoFields } from 'components/form/CustomerContactInfoFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { Form } from 'components/form/Form';

const ContactInfoForm = () => {
  return (
    <Formik initialValues={{}} onSubmit={() => {}}>
      {() => {
        return (
          <Form>
            <h1>Your contact info</h1>
            <SectionWrapper>
              <CustomerContactInfoFields />
            </SectionWrapper>
          </Form>
        );
      }}
    </Formik>
  );
};

export default ContactInfoForm;
