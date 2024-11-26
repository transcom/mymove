import React from 'react';
import { Formik } from 'formik';
import { useNavigate, useParams, generatePath } from 'react-router-dom';
import { FormGroup } from '@material-ui/core';
import classnames from 'classnames';

import SectionWrapper from 'components/Customer/SectionWrapper';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { primeSimulatorRoutes } from 'constants/routes';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import TextField from 'components/form/fields/TextField/TextField';

const PrimeUIShipmentUpdateDestinationAddressForm = ({
  initialValues,
  onSubmit,
  updateDestinationAddressSchema,
  name,
}) => {
  const { moveCodeOrID } = useParams();
  const navigate = useNavigate();

  const handleClose = () => {
    navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };

  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit} validationSchema={updateDestinationAddressSchema}>
      {({ isValid, isSubmitting, handleSubmit, errors, ...formikProps }) => (
        <Form className={classnames(formStyles.form)}>
          <FormGroup error={errors != null && Object.keys(errors).length > 0 ? 1 : 0}>
            <SectionWrapper className={formStyles.formSection}>
              <h2>Update Shipment Delivery Address</h2>
              <SectionWrapper className={formStyles.formSection}>
                <div data-testid="destination-form-details">
                  This is used to <strong>update</strong> the delivery address on an <strong>already approved</strong>{' '}
                  shipment. <br />
                  This also updates the final delivery address for destination SIT service items in the shipment.
                  <br />
                  <br />
                  This endpoint should be used for changing the delivery address of HHG & NTSR shipments.
                  <br />
                  <br />
                  The address update will be automatically approved unless it changes any of the following:
                  <br />
                  <strong>
                    - the service area <br />
                    - mileage bracket for direct delivery <br />
                    - domestic short haul to domestic line haul or vice versa <br />- SIT delivery out over 50 miles{' '}
                    <em>or </em>
                    back under 50 miles
                  </strong>
                  <br />
                  <br />
                  If any of those change, the address change will require TOO approval.
                </div>
              </SectionWrapper>
              <AddressFields name={name} locationLookup formikProps={formikProps} />
              <TextField label="Contractor Remarks" id="contractorRemarks" name="contractorRemarks" />
            </SectionWrapper>
            <WizardNavigation
              editMode
              className={formStyles.formActions}
              aria-label="Update Shipment Delivery Address"
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

export default PrimeUIShipmentUpdateDestinationAddressForm;
