import React from 'react';
import { Formik } from 'formik';
import { useNavigate, useParams, generatePath } from 'react-router-dom';
import { FormGroup } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './PrimeUIUpdateSITForms.module.scss';

import SectionWrapper from 'components/Customer/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField/TextField';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { primeSimulatorRoutes } from 'constants/routes';
import { DatePickerInput } from 'components/form/fields';
import { SERVICE_ITEM_STATUSES } from 'constants/serviceItems';

const PrimeUIUpdateOriginSITForm = ({ initialValues, onSubmit, serviceItem }) => {
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
            <div className={styles.Sit}>
              <h2 className={styles.sitHeader}>Update Origin SIT Service Item</h2>
              <SectionWrapper className={formStyles.formSection}>
                <div className={styles.sitHeader}>
                  Here you can update specific fields for an origin SIT service item. <br />
                  At this time, only the following values can be updated: <br />{' '}
                  <strong>
                    SIT Departure Date <br />
                    SIT Requested Delivery <br />
                    SIT Customer Contacted <br />
                    Update Reason
                  </strong>
                  <br />
                </div>
              </SectionWrapper>
              <SectionWrapper className={formStyles.formSection}>
                <h3 className={styles.sitHeader}>
                  {serviceItem.reServiceCode} - {serviceItem.reServiceName}
                </h3>
                <dl className={descriptionListStyles.descriptionList}>
                  <div className={descriptionListStyles.row}>
                    <dt>ID:</dt>
                    <dd>{serviceItem.id}</dd>
                  </div>
                  <div className={descriptionListStyles.row}>
                    <dt>MTO ID:</dt>
                    <dd>{serviceItem.moveTaskOrderID}</dd>
                  </div>
                  <div className={descriptionListStyles.row}>
                    <dt>Shipment ID:</dt>
                    <dd>{serviceItem.mtoShipmentID}</dd>
                  </div>
                  <div className={descriptionListStyles.row}>
                    <dt>Status:</dt>
                    <dd>{serviceItem.status}</dd>
                  </div>
                </dl>
                <div className={styles.sitDatePickerRow}>
                  <DatePickerInput name="sitDepartureDate" label="SIT Departure Date" />
                  <DatePickerInput name="sitRequestedDelivery" label="SIT Requested Delivery" />
                  <DatePickerInput name="sitCustomerContacted" label="SIT Customer Contacted" />
                </div>
                {serviceItem.status === SERVICE_ITEM_STATUSES.REJECTED && (
                  <TextField
                    display="textarea"
                    label="Update Reason"
                    data-testid="updateReason"
                    name="updateReason"
                    className={`${formStyles.remarks}`}
                    placeholder=""
                    id="updateReason"
                    maxLength={500}
                  />
                )}
              </SectionWrapper>
              <WizardNavigation
                editMode
                className={formStyles.formActions}
                aria-label="Update SIT Service Item"
                type="submit"
                onCancelClick={handleClose}
                onNextClick={handleSubmit}
              />
            </div>
          </FormGroup>
        </Form>
      )}
    </Formik>
  );
};

export default PrimeUIUpdateOriginSITForm;
