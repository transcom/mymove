import React from 'react';
import { Button, Textarea, Label, FormGroup, Radio } from '@trussworks/react-uswds'; // Tag Label
import { Formik, Field } from 'formik';

import Modal, { ModalActions, ModalClose, connectModal } from 'components/Modal/Modal'; // ModalTitle
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';

export const ShipmentAddressUpdateReviewRequestModal = ({ onClose }) => {
  return (
    <Modal>
      <ModalClose handleClick={() => onClose()} />
      <Formik>
        <Form className={formStyles.form}>
          <div>
            <div>
              <h3>Address update form</h3>
            </div>
            <FormGroup>
              <h3 style={{ fontSize: '17px' }}>Review Request</h3>
              <Label>Approve address change?</Label>
              <div data-testid="reviewSITAddressUpdateForm">
                <Field
                  as={Radio}
                  label="Yes"
                  id="acceptAddressUpdate"
                  name="sitAddressUpdate"
                  value="YES"
                  type="radio"
                />
                <Field as={Radio} label="No" id="rejectAddressUpdate" name="sitAddressUpdate" value="NO" type="radio" />
              </div>
            </FormGroup>
            <Label htmlFor="officeRemarks">Office remarks</Label>
            <p style={{ fontSize: 'small' }}>Office remarks will be sent to the contractor.</p>
            <Field as={Textarea} data-testid="officeRemarks" label="No" name="officeRemarks" id="officeRemarks" />
          </div>
          <ModalActions>
            <Button type="submit" disabled={false}>
              Save
            </Button>
            <Button type="button" secondary onClick={() => onClose()}>
              Cancel
            </Button>
          </ModalActions>
        </Form>
      </Formik>
    </Modal>
  );
};

ShipmentAddressUpdateReviewRequestModal.propTypes = {};

ShipmentAddressUpdateReviewRequestModal.defaultProps = {};

ShipmentAddressUpdateReviewRequestModal.displayName = 'ShipmentAddressUpdateReviewRequestModal';
export default connectModal(ShipmentAddressUpdateReviewRequestModal);
