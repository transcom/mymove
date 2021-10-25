import React from 'react';
import { Formik } from 'formik';
import PropTypes from 'prop-types';
import { FormGroup } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import classnames from 'classnames';

import SectionWrapper from '../../../components/Customer/SectionWrapper';
import { ResidentialAddressShape } from '../../../types/address';
import { AddressFields } from '../../../components/form/AddressFields/AddressFields';

import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { requiredAddressSchema } from 'utils/validation';

const PrimeUIShipmentUpdateAddressForm = ({
  initialValues,
  addressLocation,
  onSubmit,
  updateShipmentAddressSchema,
}) => (
  <Formik initialValues={initialValues} onSubmit={onSubmit} validationSchema={updateShipmentAddressSchema}>
    {({ isValid, isSubmitting, errors }) => (
      /* <Form className={classnames(styles.CreatePaymentRequestForm, formStyles.form)}> */
      <Form className={classnames(formStyles.form)}>
        <FormGroup error={errors != null && Object.keys(errors).length > 0 ? 1 : 0}>
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
    addressID: PropTypes.string.isRequired,
    eTag: PropTypes.string.isRequired,
  }).isRequired,
  onSubmit: PropTypes.func.isRequired,
  updateShipmentAddressSchema: PropTypes.shape({
    address: requiredAddressSchema,
  }).isRequired,
  addressLocation: PropTypes.string.isRequired,
};

export default PrimeUIShipmentUpdateAddressForm;
