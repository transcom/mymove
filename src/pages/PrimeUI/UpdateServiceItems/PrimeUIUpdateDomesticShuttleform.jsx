import React from 'react';
import { Formik } from 'formik';
import { useNavigate, useParams, generatePath } from 'react-router-dom';
import { FormGroup } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './PrimeUIUpdateShuttleForms.module.scss';

import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField/TextField';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { primeSimulatorRoutes } from 'constants/routes';
import { SERVICE_ITEM_STATUSES } from 'constants/serviceItems';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { CheckboxField } from 'components/form/fields';
import { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

const PrimeUIUpdateDomesticShuttleForm = ({ onUpdateServiceItem, serviceItem }) => {
  const { moveCodeOrID } = useParams();
  const navigate = useNavigate();

  const handleClose = () => {
    navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };

  const initialValues = {
    mtoServiceItemID: serviceItem.id,
    reServiceCode: serviceItem.reServiceCode,
    eTag: serviceItem.eTag,
  };

  const onSubmit = (values) => {
    const {
      eTag,
      mtoServiceItemID,
      estimatedWeight,
      actualWeight,
      updateReason,
      reServiceCode,
      requestApprovalsRequestedStatus,
    } = values;

    const body = {
      reServiceCode,
      modelType: 'UpdateMTOServiceItemShuttle',
      actualWeight: Number.parseInt(actualWeight, 10),
      estimatedWeight: Number.parseInt(estimatedWeight, 10),
      updateReason,
      requestApprovalsRequestedStatus,
    };

    onUpdateServiceItem({ mtoServiceItemID, eTag, body });
  };

  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit}>
      {({ handleSubmit, setFieldValue }) => (
        <Form className={classnames(formStyles.form)}>
          <FormGroup>
            <div className={styles.Shuttle}>
              <h2 className={styles.shuttleHeader}>Update Domestic Shuttle Service Item</h2>
              <SectionWrapper className={formStyles.formSection}>
                <div className={styles.shuttleHeader}>
                  Here you can update specific fields for an domestic shuttle service item. <br />
                  At this time, only the following values can be updated: <br />{' '}
                  <strong>
                    Estimated Weight <br />
                    Actual Weight <br />
                    Update Reason
                  </strong>
                  <br />
                </div>
              </SectionWrapper>
              <SectionWrapper className={formStyles.formSection}>
                <h3 className={styles.shuttleHeader}>
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
                  {serviceItem.status === SERVICE_ITEM_STATUSES.REJECTED && (
                    <div className={descriptionListStyles.row}>
                      <dt>Rejection Reason:</dt>
                      <dd>{serviceItem.rejectionReason}</dd>
                    </div>
                  )}
                  <div className={descriptionListStyles.row}>
                    <dt>Estimated Weight:</dt>
                    <dd>{serviceItem.estimatedWeight}</dd>
                  </div>
                  <div className={descriptionListStyles.row}>
                    <dt>Actual Weight:</dt>
                    <dd>{serviceItem.actualWeight}</dd>
                  </div>
                </dl>
                <MaskedTextField
                  data-testid="estimatedWeightInput"
                  name="estimatedWeight"
                  label="Estimated weight (lbs)"
                  id="estimatedWeightInput"
                  mask={Number}
                  scale={0}
                  thousandsSeparator=","
                  lazy={false} // immediate masking evaluation
                  onChange={(e) => {
                    setFieldValue('estimatedWeight', e.target.value);
                  }}
                />
                <MaskedTextField
                  data-testid="actualWeightInput"
                  name="actualWeight"
                  label="Actual Weight"
                  id="actualWeight"
                  mask={Number}
                  scale={0}
                  thousandsSeparator=","
                  lazy={false} // immediate masking evaluation
                  onChange={(e) => {
                    setFieldValue('actualWeight', e.target.value);
                  }}
                />
                {serviceItem.status === SERVICE_ITEM_STATUSES.REJECTED && (
                  <>
                    {requiredAsteriskMessage}
                    <TextField
                      display="textarea"
                      label="Update Reason"
                      data-testid="updateReason"
                      name="updateReason"
                      className={`${formStyles.remarks}`}
                      placeholder=""
                      id="updateReason"
                      maxLength={500}
                      showRequiredAsterisk
                      required
                    />
                  </>
                )}
                {serviceItem.status === SERVICE_ITEM_STATUSES.REJECTED && (
                  <CheckboxField
                    data-testid="requestApprovalsRequestedStatus"
                    label="Request Approval"
                    name="requestApprovalsRequestedStatus"
                    id="requestApprovalsRequestedStatus"
                  />
                )}
              </SectionWrapper>
              <WizardNavigation
                editMode
                className={formStyles.formActions}
                aria-label="Update Domestic Shuttle Service Item"
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

export default PrimeUIUpdateDomesticShuttleForm;
