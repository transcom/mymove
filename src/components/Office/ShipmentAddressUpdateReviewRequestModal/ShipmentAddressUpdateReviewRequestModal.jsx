import React from 'react';
import { Button, Textarea, Label, FormGroup, Radio } from '@trussworks/react-uswds'; // Tag Label
import { Formik, Field } from 'formik';
import * as Yup from 'yup';

import styles from './ShipmentAddressUpdateReviewRequestModal.module.scss';

import Modal, { ModalActions, ModalClose, ModalTitle, connectModal } from 'components/Modal/Modal'; // ModalTitle
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import AddressUpdatePreview from 'components/Office/AddressUpdatePreview/AddressUpdatePreview';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';

const formSchema = Yup.object().shape({
  addressUpdate: Yup.string().required('Required'),
  officeRemarks: Yup.string().required('Required'),
});

export const ShipmentAddressUpdateReviewRequestModal = ({ deliveryAddressUpdate, shipmentType, onClose }) => {
  return (
    <Modal>
      <ModalClose handleClick={() => onClose()} />
      <ModalTitle>
        <ShipmentTag shipmentType={shipmentType} />
        <h2 className={styles.modalTitle}>Review request</h2>
      </ModalTitle>
      <Formik
        initialValues={{ addressUpdate: '', officeRemarks: '' }}
        onSubmit={() => {}}
        validateOnMount
        validationSchema={formSchema}
      >
        {({ isValid }) => {
          return (
            <Form className={formStyles.form}>
              <div className={styles.modalbody}>
                <AddressUpdatePreview deliveryAddressUpdate={deliveryAddressUpdate} shipmentType={shipmentType} />
                <FormGroup>
                  <h4>Review Request</h4>
                  <Label>Approve address change?</Label>
                  <div>
                    <Field
                      as={Radio}
                      label="Yes"
                      id="acceptAddressUpdate"
                      name="addressUpdate"
                      value="YES"
                      type="radio"
                    />
                    <Field
                      as={Radio}
                      label="No"
                      id="rejectAddressUpdate"
                      name="addressUpdate"
                      value="NO"
                      type="radio"
                    />
                  </div>
                </FormGroup>
                <Label htmlFor="officeRemarks">Office remarks</Label>
                <p style={{ fontSize: 'small' }}>Office remarks will be sent to the contractor.</p>
                <Field as={Textarea} data-testid="officeRemarks" label="No" name="officeRemarks" id="officeRemarks" />
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

ShipmentAddressUpdateReviewRequestModal.propTypes = {};

ShipmentAddressUpdateReviewRequestModal.defaultProps = {};

ShipmentAddressUpdateReviewRequestModal.displayName = 'ShipmentAddressUpdateReviewRequestModal';
export default connectModal(ShipmentAddressUpdateReviewRequestModal);
