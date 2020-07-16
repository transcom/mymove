import React from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { Formik } from 'formik';

import { Form } from '../form/Form';
import { DatePickerInput } from '../form/fields';
import { AddressFields } from '../form/AddressFields/AddressFields';
import { ContactInfoFields } from '../form/ContactInfoFields/ContactInfoFields';
import { WizardPage } from '../../shared/WizardPage';

import { TextInput } from 'components/form/fields/TextInput';

// eslint-disable-next-line
export const HHGDetailsForm = ({ initialValues, pageKey, pages }) => {
  return (
    <Formik
      initialValues={initialValues}
      onSubmit={(values) => {
        console.log('this is values', values);
      }}
    >
      {({ handleSubmit }) => (
        <WizardPage pageKey={pageKey} pageList={pages} handleSubmit={handleSubmit}>
          <Form>
            <DatePickerInput name="requestedPickupDate" label="Requested pickup date" id="requested-pickup-date" />
            <AddressFields addressType="pickup" legend="Pickup location" />
            <ContactInfoFields contactType="releasing" legend="Releasing agent" />
            <DatePickerInput
              name="requestedDeliveryDate"
              label="Requested delivery date"
              id="requested-delivery-date"
            />
            <AddressFields addressType="delivery" legend="Delivery location" />
            <ContactInfoFields contactType="receiving" legend="Receiving agent" />
            <TextInput name="remarks" label="Remarks" hint=" (optional)" id="requested-delivery-date" />
          </Form>
        </WizardPage>
      )}
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
