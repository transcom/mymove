import React from 'react';
import { Formik } from 'formik';
import { useHistory, useParams } from 'react-router-dom-old';
import { generatePath } from 'react-router';
import PropTypes from 'prop-types';
import { FormGroup } from '@material-ui/core';
import classnames from 'classnames';

import SectionWrapper from 'components/Customer/SectionWrapper';
import { ResidentialAddressShape } from 'types/address';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { primeSimulatorRoutes } from 'constants/routes';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';

const PrimeUIShipmentUpdateAddressForm = ({
  initialValues,
  addressLocation,
  onSubmit,
  updateShipmentAddressSchema,
  name,
}) => {
  const { moveCodeOrID } = useParams();
  const history = useHistory();

  const handleClose = () => {
    history.push(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };

  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit} validationSchema={updateShipmentAddressSchema}>
      {({ isValid, isSubmitting, handleSubmit, errors }) => (
        <Form className={classnames(formStyles.form)}>
          <FormGroup error={errors != null && Object.keys(errors).length > 0 ? 1 : 0}>
            <SectionWrapper className={formStyles.formSection}>
              <h2>{addressLocation}</h2>
              <AddressFields name={name} />
            </SectionWrapper>
            <WizardNavigation
              editMode
              className={formStyles.formActions}
              aria-label="Update Shipment Address"
              type="submit"
              disableNext={isSubmitting || !isValid}
              onCancelClick={handleClose}
              onNextClick={handleSubmit}
            />
          </FormGroup>
        </Form>
      )}
    </Formik>
  );
};

PrimeUIShipmentUpdateAddressForm.propTypes = {
  initialValues: PropTypes.shape({
    pickupAddress: PropTypes.shape({
      address: ResidentialAddressShape,
    }),
    destinationAddress: PropTypes.shape({
      address: ResidentialAddressShape,
    }),
    addressID: PropTypes.string,
    eTag: PropTypes.string,
  }).isRequired,
  onSubmit: PropTypes.func.isRequired,
  updateShipmentAddressSchema: PropTypes.shape({
    address: ResidentialAddressShape,
    addressID: PropTypes.string,
    eTag: PropTypes.string,
  }).isRequired,
  addressLocation: PropTypes.oneOf(['Pickup address', 'Destination address']).isRequired,
  name: PropTypes.string.isRequired,
};

export default PrimeUIShipmentUpdateAddressForm;
