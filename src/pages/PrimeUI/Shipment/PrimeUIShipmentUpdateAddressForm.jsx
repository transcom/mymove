import React from 'react';
import { Formik } from 'formik';
import { useHistory, useParams } from 'react-router-dom';
import { generatePath } from 'react-router';
import PropTypes from 'prop-types';
import { FormGroup } from '@material-ui/core';
import classnames from 'classnames';

import SectionWrapper from '../../../components/Customer/SectionWrapper';
import { ResidentialAddressShape } from '../../../types/address';
import { AddressFields } from '../../../components/form/AddressFields/AddressFields';
import { primeSimulatorRoutes } from '../../../constants/routes';

import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { requiredAddressSchema } from 'utils/validation';

const PrimeUIShipmentUpdateAddressForm = ({
  initialValues,
  addressLocation,
  onSubmit,
  updateShipmentAddressSchema,
}) => {
  const { moveCodeOrID } = useParams();
  const history = useHistory();

  const handleClose = () => {
    history.push(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };

  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit} validationSchema={updateShipmentAddressSchema}>
      {({ isValid, isSubmitting, handleSubmit, errors }) => (
        /* <Form className={classnames(styles.CreatePaymentRequestForm, formStyles.form)}> */
        <Form className={classnames(formStyles.form)}>
          <FormGroup error={errors != null && Object.keys(errors).length > 0 ? 1 : 0}>
            <SectionWrapper className={formStyles.formSection}>
              <h2>{addressLocation}</h2>
              <AddressFields name="address" />
            </SectionWrapper>
            <WizardNavigation
              editMode
              className={formStyles.formActions}
              aria-label="Update Shipment Address"
              type="submit"
              disabled={isSubmitting || !isValid}
              onCancelClick={handleClose}
              onNextClick={handleSubmit}
            >
              Update
            </WizardNavigation>
          </FormGroup>
        </Form>
      )}
    </Formik>
  );
};

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
