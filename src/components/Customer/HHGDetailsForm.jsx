// HHG details form storybook component
import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
// import { Button } from '@trussworks/react-uswds';

import { Form } from '../form/Form';
import { DatePickerInput, TextInput } from '../form/fields';

// eslint-disable-next-line
export const HHGDetailsForm = ({ initialValues }) => {
  return (
    <Formik initialValues={{ remarks: '' }}>
      <Form>
        <DatePickerInput name="requestedPickupDate" label="Requested pickup date" />
        <DatePickerInput name="requestedDeliveryDate" label="Requested delivery date" />
        <TextInput name="remarks" label="Remarks" />
      </Form>
    </Formik>
  );
};

HHGDetailsForm.propTypes = {
  initialValues: PropTypes.shape({
    remarks: PropTypes.string,
    requestedPickupDate: PropTypes.string,
    requestedDeliveryDate: PropTypes.string,
  }),
};

HHGDetailsForm.defaultProps = {
  initialValues: {},
};

export default HHGDetailsForm;
