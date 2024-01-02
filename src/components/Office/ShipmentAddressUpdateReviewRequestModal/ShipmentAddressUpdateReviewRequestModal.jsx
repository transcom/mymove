import React from 'react';
import { Alert, Button, Textarea, Label, FormGroup, Radio } from '@trussworks/react-uswds'; // Tag Label
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import * as PropTypes from 'prop-types';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './ShipmentAddressUpdateReviewRequestModal.module.scss';

import Modal, { ModalActions, ModalClose, ModalTitle, connectModal } from 'components/Modal/Modal'; // ModalTitle
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import AddressUpdatePreview from 'components/Office/AddressUpdatePreview/AddressUpdatePreview';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { ShipmentShape } from 'types';
import Fieldset from 'shared/Fieldset';
import { ADDRESS_UPDATE_STATUS } from 'constants/shipments';

const formSchema = Yup.object().shape({
  addressUpdateReviewStatus: Yup.string().required('Required'),
  officeRemarks: Yup.string().required('Required'),
});

export const ShipmentAddressUpdateReviewRequestModal = ({
  onSubmit,
  shipment,
  errorMessage,
  setErrorMessage,
  onClose,
}) => {
  const handleSubmit = async (values, { setSubmitting }) => {
    const { addressUpdateReviewStatus, officeRemarks } = values;

    await onSubmit(shipment.id, shipment.eTag, addressUpdateReviewStatus, officeRemarks);

    setSubmitting(false);
  };

  const errorMessageAlertControl = (
    <Button type="button" onClick={() => setErrorMessage(null)} unstyled>
      <FontAwesomeIcon icon="times" style={styles.alertClose} />
    </Button>
  );

  // checking to see if the shipment contains destination SIT service items
  // storing them in an array so we can have the count and display to the TOO
  // if there is none, we will make it an empty array
  const destSitServiceItems = (shipment.mtoServiceItems ?? []).filter((s) =>
    ['DDDSIT', 'DDASIT', 'DDFSIT', 'DDSFSC'].includes(s.reServiceCode),
  );

  return (
    <Modal>
      <ModalClose handleClick={() => onClose()} />
      <ModalTitle>
        <ShipmentTag shipmentType={shipment.shipmentType} />
        <h2 className={styles.modalTitle}>Review request</h2>
        {errorMessage && (
          <Alert type="error" role="alert" cta={errorMessageAlertControl}>
            {errorMessage}
          </Alert>
        )}
      </ModalTitle>
      <Formik
        initialValues={{ addressUpdateReviewStatus: '', officeRemarks: '' }}
        onSubmit={handleSubmit}
        validateOnMount
        validationSchema={formSchema}
      >
        {({ isValid }) => {
          return (
            <Form className={formStyles.form}>
              <div className={styles.modalbody}>
                <AddressUpdatePreview
                  deliveryAddressUpdate={shipment.deliveryAddressUpdate}
                  shipmentType={shipment.shipmentType}
                  destSitServiceItems={destSitServiceItems}
                />
                <FormGroup className={styles.formGroup}>
                  <h4>Review Request</h4>
                  <Fieldset>
                    <legend className={classnames('usa-label', styles.approveLabel)}>Approve address change?</legend>
                    <Field
                      as={Radio}
                      label="Yes"
                      id="acceptAddressUpdate"
                      name="addressUpdateReviewStatus"
                      value={ADDRESS_UPDATE_STATUS.APPROVED}
                      type="radio"
                    />
                    <Field
                      as={Radio}
                      label="No"
                      id="rejectAddressUpdate"
                      name="addressUpdateReviewStatus"
                      value={ADDRESS_UPDATE_STATUS.REJECTED}
                      type="radio"
                    />
                  </Fieldset>
                </FormGroup>
                <Label htmlFor="officeRemarks">Office remarks</Label>
                <p className={styles.subLabel}>Office remarks will be sent to the contractor.</p>
                <Field
                  as={Textarea}
                  data-testid="officeRemarks"
                  label="No"
                  name="officeRemarks"
                  id="officeRemarks"
                  className={styles.officeRemarks}
                />
              </div>
              <ModalActions>
                <Button type="submit" disabled={!isValid}>
                  Save
                </Button>
                <Button type="button" secondary onClick={onClose}>
                  Cancel
                </Button>
              </ModalActions>
            </Form>
          );
        }}
      </Formik>
    </Modal>
  );
};

ShipmentAddressUpdateReviewRequestModal.propTypes = {
  shipment: ShipmentShape.isRequired,
  onSubmit: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
  errorMessage: PropTypes.node,
  setErrorMessage: PropTypes.func,
};

ShipmentAddressUpdateReviewRequestModal.defaultProps = {
  errorMessage: null,
  setErrorMessage: undefined,
};

ShipmentAddressUpdateReviewRequestModal.displayName = 'ShipmentAddressUpdateReviewRequestModal';
export default connectModal(ShipmentAddressUpdateReviewRequestModal);
