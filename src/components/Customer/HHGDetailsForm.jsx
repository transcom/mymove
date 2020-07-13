import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';

import { Form } from '../form/Form';
import { DatePickerInput, TextInput } from '../form/fields';
import { AddressFields } from '../form/AddressFields/AddressFields';
import { ContactInfoFields } from '../form/ContactInfoFields/ContactInfoFields';
import { WizardPage } from '../../shared/WizardPage';

// eslint-disable-next-line
export const HHGDetailsForm = ({ initialValues, pageKey, pages }) => {
  return (
    <Formik initialValues={{ remarks: '' }}>
      <WizardPage pageKey={pageKey} pageList={pages}>
        <Form>
          <DatePickerInput name="requestedPickupDate" label="Requested pickup date" />
          <AddressFields initialValues={initialValues.pickupLocation} legend="Pickup location" />
          <ContactInfoFields legend="Releasing agent" />
          <DatePickerInput name="requestedDeliveryDate" label="Requested delivery date" />
          <AddressFields initialValues={initialValues.deliveryLocation} legend="Delivery location" />
          <ContactInfoFields legend="Receiving agent" />
          <TextInput name="remarks" label="Remarks" />
        </Form>
      </WizardPage>
    </Formik>
  );
};

HHGDetailsForm.propTypes = {
  initialValues: PropTypes.shape({
    requestedPickupDate: PropTypes.string,
    pickupLocation: PropTypes.shape({
      mailingAddress1: PropTypes.string,
      mailingAddress2: PropTypes.string,
      city: PropTypes.string,
      state: PropTypes.string,
      zip: PropTypes.string,
    }),
    requestedDeliveryDate: PropTypes.string,
    deliveryLocation: PropTypes.shape({
      mailingAddress1: PropTypes.string,
      mailingAddress2: PropTypes.string,
      city: PropTypes.string,
      state: PropTypes.string,
      zip: PropTypes.string,
    }),
    remarks: PropTypes.string,
  }),
};

HHGDetailsForm.defaultProps = {
  initialValues: {},
};

export default HHGDetailsForm;
