import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';

import { Form } from '../form/Form';
import { DatePickerInput, TextInput } from '../form/fields';
import { AddressFields } from '../form/AddressFields';

// eslint-disable-next-line
export const HHGDetailsForm = ({ initialValues }) => {
  return (
    <Formik initialValues={{ remarks: '' }}>
      <Form>
        <DatePickerInput name="requestedPickupDate" label="Requested pickup date" />
        <AddressFields legend="Pickup location" />
        <DatePickerInput name="requestedDeliveryDate" label="Requested delivery date" />
        <AddressFields legend="Delivery location" />
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
