import React from 'react';
import { Formik } from 'formik';
import { useNavigate, useParams, generatePath } from 'react-router-dom';
import { FormGroup } from '@trussworks/react-uswds';
import classnames from 'classnames';

import SectionWrapper from 'components/Customer/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import { Form } from 'components/form/Form';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { primeSimulatorRoutes } from 'constants/routes';
import { DatePickerInput } from 'components/form/fields';

const PrimeUIUpdateOriginSITForm = ({ initialValues, onSubmit }) => {
  const { moveCodeOrID } = useParams();
  const navigate = useNavigate();

  const handleClose = () => {
    navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };

  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit}>
      {({ handleSubmit }) => (
        <Form className={classnames(formStyles.form)}>
          <FormGroup>
            <h2 style={{ textAlign: 'center' }}>Update Origin SIT Service Item</h2>
            <SectionWrapper className={formStyles.formSection}>
              Here you can update specific fields for an origin SIT service item. <br />
              At this time, only the <strong>SIT Departure Date</strong>, <strong>SIT Requested Delivery</strong>, and{' '}
              <strong>SIT Customer Contacted</strong> fields can be updated.
            </SectionWrapper>
            <SectionWrapper className={formStyles.formSection}>
              <h2 style={{ textAlign: 'center' }}>Update Origin SIT Service Item</h2>
              <div style={{ display: 'flex', justifyContent: 'space-around' }}>
                <DatePickerInput name="sitDepartureDate" label="SIT Departure Date" />
                <DatePickerInput name="sitRequestedDelivery" label="SIT Requested Delivery" />
                <DatePickerInput name="sitCustomerContacted" label="SIT Customer Contacted" />
              </div>
            </SectionWrapper>
            <WizardNavigation
              editMode
              className={formStyles.formActions}
              aria-label="Update SIT Service Item"
              type="submit"
              onCancelClick={handleClose}
              onNextClick={handleSubmit}
            />
          </FormGroup>
        </Form>
      )}
    </Formik>
  );
};

export default PrimeUIUpdateOriginSITForm;
