import React, { useState } from 'react';
import { Formik } from 'formik';
import PropTypes from 'prop-types';

import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { requiredAddressSchema } from 'utils/validation';
import { FormGroup } from '@material-ui/core';
import SectionWrapper from '../../../components/Customer/SectionWrapper';
import { ErrorMessage } from '../../../components/form';
import { ResidentialAddressShape } from '../../../types/address';
import AddressFields from '../../../components/form/AddressFields/AddressFields';
import Button from '@material-ui/core/Button';

const PrimeUIShipmentUpdateAddressForm = ({
  initialValues,
  addressLocation,
  onSubmit,
  updateShipmentAddressSchema,
}) => (
  <Formik initialValues={initialValues} onSubmit={onSubmit} validationSchema={updateShipmentAddressSchema}>
    {({ isValid, isSubmitting, errors, values, setValues }) => (
      /* <Form className={classnames(styles.CreatePaymentRequestForm, formStyles.form)}> */
      <Form className={classnames(formStyles.form)}>
        <FormGroup error={errors != null && Object.keys(errors).length > 0}>
          {errors != null /* && errors.serviceItems */ && (
            <ErrorMessage display>At least 1 service item must be added when creating a payment request</ErrorMessage>
          )}
          <SectionWrapper className={formStyles.formSection}>
            <h2>{addressLocation}</h2>
            <AddressFields name="address" />
          </SectionWrapper>
          <Button aria-label="Update Shipment Address" type="submit" disabled={isSubmitting || !isValid}>
            Update
          </Button>
        </FormGroup>
      </Form>
    )}
  </Formik>
);

PrimeUIShipmentUpdateAddressForm.propTypes = {
  initialValues: PropTypes.shape({
    address: ResidentialAddressShape,
  }).isRequired,
  onSubmit: PropTypes.func.isRequired,
  updateShipmentAddressSchema: PropTypes.shape({
    address: requiredAddressSchema,
  }).isRequired,
  addressLocation: PropTypes.string.isRequired,
};

export default PrimeUIShipmentUpdateAddressForm;
