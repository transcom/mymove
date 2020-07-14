import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';

import { Form } from '../form/Form';
import { DatePickerInput, TextInput } from '../form/fields';
import { AddressFields } from '../form/AddressFields/AddressFields';
import { ContactInfoFields } from '../form/ContactInfoFields/ContactInfoFields';
import { WizardPage } from '../../shared/WizardPage';

export const HHGDetailsForm = ({ initialValues, pageKey, pageList }) => {
  return (
    <Formik initialValues={initialValues}>
      <WizardPage pageKey={pageKey} pageList={pageList} handleSubmit={() => {}}>
        <Form>
          <DatePickerInput name="requestedPickupDate" label="Requested pickup date" id="requested-pickup-date" />
          <AddressFields initialValues={initialValues.pickupLocation} legend="Pickup location" />
          <ContactInfoFields legend="Releasing agent" />
          <DatePickerInput name="requestedDeliveryDate" label="Requested delivery date" id="requested-delivery-date" />
          <AddressFields initialValues={initialValues.deliveryLocation} legend="Delivery location" />
          <ContactInfoFields legend="Receiving agent" />
          <TextInput name="remarks" label="Remarks" id="requested-delivery-date" />
        </Form>
      </WizardPage>
    </Formik>
  );
};

HHGDetailsForm.propTypes = {
  pageKey: PropTypes.string.isRequired,
  pageList: PropTypes.arrayOf(PropTypes.string).isRequired,
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
