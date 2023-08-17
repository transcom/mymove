import React from 'react';
import { Formik } from 'formik';
import { useNavigate, useParams, generatePath } from 'react-router-dom';
import PropTypes from 'prop-types';
import { FormGroup } from '@trussworks/react-uswds';
import classnames from 'classnames';

import SectionWrapper from 'components/Customer/SectionWrapper';
import AddressFields from 'components/form/AddressFields/AddressFields';
import TextField from 'components/form/fields/TextField/TextField';
import { ResidentialAddressShape } from 'types/address';
import formStyles from 'styles/form.module.scss';
import { Form } from 'components/form/Form';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { primeSimulatorRoutes } from 'constants/routes';

const PrimeUIRequestSITDestAddressChangeForm = ({ name, initialValues, onSubmit, destAddressChangeRequestSchema }) => {
  const { moveCodeOrID } = useParams();
  const navigate = useNavigate();

  const handleClose = () => {
    navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };

  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit} validationSchema={destAddressChangeRequestSchema}>
      {({ isValid, isSubmitting, handleSubmit, errors }) => (
        <Form className={classnames(formStyles.form)}>
          <FormGroup error={errors != null && Object.keys(errors).length > 0 ? 1 : 0}>
            <SectionWrapper className={formStyles.formSection}>
              <h2>Request Destination SIT Address Change </h2>
              <AddressFields name={name} />
              <TextField name="contractorRemarks" id="contractorRemarks" label="Contractor Remarks" />
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

PrimeUIRequestSITDestAddressChangeForm.propTypes = {
  initialValues: PropTypes.shape({
    address: PropTypes.shape({ ResidentialAddressShape }),
    contractorRemarks: PropTypes.string,
    mtoServiceItemID: PropTypes.string,
  }).isRequired,
  destAddressChangeRequestSchema: PropTypes.shape({
    address: ResidentialAddressShape,
    contractorRemarks: PropTypes.string,
    mtoServiceItemID: PropTypes.string,
  }).isRequired,
  name: PropTypes.string.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

export default PrimeUIRequestSITDestAddressChangeForm;
