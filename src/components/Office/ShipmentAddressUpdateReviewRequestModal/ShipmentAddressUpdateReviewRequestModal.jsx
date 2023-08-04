import React from 'react';
import { Button, Textarea, Label, FormGroup, Radio } from '@trussworks/react-uswds'; // Tag Label
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import * as PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './ShipmentAddressUpdateReviewRequestModal.module.scss';

import Modal, { ModalActions, ModalClose, ModalTitle, connectModal } from 'components/Modal/Modal'; // ModalTitle
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import AddressUpdatePreview from 'components/Office/AddressUpdatePreview/AddressUpdatePreview';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { ShipmentAddressUpdateShape } from 'types';
import Fieldset from 'shared/Fieldset';
import { ADDRESS_UPDATE_STATUS } from 'constants/shipments';

const formSchema = Yup.object().shape({
  addressUpdate: Yup.string().required('Required'),
  officeRemarks: Yup.string().required('Required'),
});

export const ShipmentAddressUpdateReviewRequestModal = ({ onSubmit, deliveryAddressUpdate, shipmentType, onClose }) => {
  return (
    <Modal>
      <ModalClose handleClick={() => onClose()} />
      <ModalTitle>
        <ShipmentTag shipmentType={shipmentType} />
        <h2 className={styles.modalTitle}>Review request</h2>
        {/* TODO: Error alert  <Alert>
          { errorMessage}
        </Alert> */}
      </ModalTitle>
      <Formik
        initialValues={{ addressUpdateReviewStatus: '', officeRemarks: '' }}
        onSubmit={onSubmit}
        validateOnMount
        validationSchema={formSchema}
      >
        {({ isValid }) => {
          return (
            <Form className={formStyles.form}>
              <div className={styles.modalbody}>
                <AddressUpdatePreview deliveryAddressUpdate={deliveryAddressUpdate} shipmentType={shipmentType} />
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
  deliveryAddressUpdate: ShipmentAddressUpdateShape.isRequired,
  shipmentType: PropTypes.string.isRequired,
  onClose: PropTypes.func.isRequired,
};

ShipmentAddressUpdateReviewRequestModal.displayName = 'ShipmentAddressUpdateReviewRequestModal';
export default connectModal(ShipmentAddressUpdateReviewRequestModal);
