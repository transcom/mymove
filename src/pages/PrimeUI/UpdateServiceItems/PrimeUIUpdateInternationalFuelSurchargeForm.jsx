import React from 'react';
import { Formik } from 'formik';
import * as Yup from 'yup';
import { useNavigate, useParams, generatePath } from 'react-router-dom';
import { FormGroup } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './PrimeUIUpdateInternationalFuelSurcharge.module.scss';

import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField/TextField';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { primeSimulatorRoutes } from 'constants/routes';
import { SERVICE_ITEM_STATUSES } from 'constants/serviceItems';
import { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

const PrimeUIUpdateInternationalFuelSurchargeForm = ({ onUpdateServiceItem, moveTaskOrder, mtoServiceItemId }) => {
  const { moveCodeOrID } = useParams();
  const navigate = useNavigate();

  const handleClose = () => {
    navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };

  const serviceItem = moveTaskOrder?.mtoServiceItems.find((s) => s?.id === mtoServiceItemId);
  const mtoShipment = moveTaskOrder?.mtoShipments.find((s) => s.id === serviceItem.mtoShipmentID);
  let port;
  if (mtoShipment.portOfEmbarkation) {
    port = mtoShipment.portOfEmbarkation;
  } else if (mtoShipment.portOfDebarkation) {
    port = mtoShipment.portOfDebarkation;
  } else {
    port = null;
  }
  const initialValues = {
    mtoServiceItemID: serviceItem.id,
    reServiceCode: serviceItem.reServiceCode,
    eTag: serviceItem.eTag,
    portCode: port?.portCode,
  };

  const onSubmit = (values) => {
    const { eTag, mtoServiceItemID, portCode, reServiceCode } = values;

    const body = {
      portCode,
      reServiceCode,
      modelType: 'UpdateMTOServiceItemInternationalPortFSC',
    };

    onUpdateServiceItem({ mtoServiceItemID, eTag, body });
  };

  return (
    <Formik
      initialValues={initialValues}
      onSubmit={onSubmit}
      validationSchema={Yup.object({
        portCode: Yup.string()
          .required('Required')
          .min(3, 'Port Code must be 3-4 characters.')
          .max(4, 'Port Code must be 3-4 characters.'),
      })}
    >
      {({ handleSubmit, setFieldValue }) => (
        <Form className={classnames(formStyles.form)}>
          <FormGroup>
            <div className={styles.IntlFsc}>
              <h2 className={styles.intlFscHeader}>Update International Fuel Surcharge Service Item</h2>
              <SectionWrapper className={formStyles.formSection}>
                <div className={styles.intlFscHeader}>
                  Here you can update specific fields for an International Fuel Surcharge service item. <br />
                  At this time, only the following values can be updated: <br />{' '}
                  <strong>
                    Port Code <br />
                  </strong>
                </div>
              </SectionWrapper>
              <SectionWrapper className={formStyles.formSection}>
                <h3 className={styles.intlFscHeader}>
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
                  <div className={descriptionListStyles.row}>
                    <dt>Port:</dt>
                    <dd>{port && port.portName}</dd>
                  </div>
                  <div className={descriptionListStyles.row}>
                    <dt>Port Location:</dt>
                    <dd>
                      {port && port.city}
                      {port && port.city && ','} {port && port.state} {port && port.zip}
                    </dd>
                  </div>
                </dl>
                {requiredAsteriskMessage}
                <TextField
                  data-testid="portCode"
                  name="portCode"
                  label="Port Code"
                  id="portCode"
                  required
                  labelHint="Required"
                  maxLength="4"
                  isDisabled={serviceItem.status !== SERVICE_ITEM_STATUSES.APPROVED}
                  onChange={(e) => {
                    setFieldValue('portCode', e.target.value.toUpperCase());
                  }}
                  showRequiredAsterisk
                />
              </SectionWrapper>
              <div className={formStyles.formActions}>
                <WizardNavigation
                  editMode
                  aria-label="Update International Fuel Surcharge Service Item"
                  type="submit"
                  onCancelClick={handleClose}
                  onNextClick={handleSubmit}
                />
              </div>
            </div>
          </FormGroup>
        </Form>
      )}
    </Formik>
  );
};

export default PrimeUIUpdateInternationalFuelSurchargeForm;
