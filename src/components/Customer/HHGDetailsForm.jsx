// HHG details form storybook component
import React from 'react';
// import PropTypes from 'prop-types';
import { Formik } from 'formik';
// import { Button } from '@trussworks/react-uswds';

import { Form } from '../form/Form';
import { TextInput } from '../form/fields';

export const HHGDetailsForm = () => {
  return (
    <Formik>
      <Form>
        <TextInput name="name" label="label" />
      </Form>
    </Formik>
  );
};

// HHGDetailsForm.propTypes = {
//   initialValues: PropTypes.shape({
//     name: PropTypes.string,
//   }),
// };

export default HHGDetailsForm;
