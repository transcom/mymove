import React from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import PropTypes from 'prop-types';

import { CustomerContactInfoFields } from 'components/form/CustomerContactInfoFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { Form } from 'components/form/Form';

const ContactInfoForm = (initialValues, onSubmit) => {
  const validationSchema = Yup.object().shape({
    telephone: Yup.string(),
    secondary_phone: Yup.string(),
    personal_email: Yup.string(),
    phone_is_preferred: Yup.bool(),
    email_is_preferred: Yup.bool(),
  });
  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit} validationSchema={validationSchema} validateOnMount>
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

ContactInfoForm.propTypes = {
  initialValues: PropTypes.shape({
    telephone: PropTypes.string,
    secondary_phone: PropTypes.string,
    personal_email: PropTypes.string,
    phone_is_preferred: PropTypes.bool,
    email_is_preferred: PropTypes.bool,
  }).isRequired,
};

export default ContactInfoForm;
