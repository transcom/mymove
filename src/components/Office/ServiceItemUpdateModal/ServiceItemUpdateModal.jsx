import React from 'react';
import PropTypes from 'prop-types';
import { Button, Label, Textarea, Tag } from '@trussworks/react-uswds';
import { Formik, Field } from 'formik';

import Modal, { ModalActions, ModalClose, ModalTitle, connectModal } from 'components/Modal/Modal';
import { Form } from 'components/form';

// const SitAddressInfo = () => {
//   return <div>Tersting</div>;
// };

/**
 * @description This componment is thee modal used for when a TOO edits the address for a Service item
 * or reviews a service item request from a the prime.
 */
export const ServiceItemUpdateModal = ({ onSave, closeModal, title }) => {
  const initialValues = {
    officeRemarks: '',
    newAddress: {
      streetAddress1: '',
      streetAddress2: '',
      city: '',
      state: '',
      postalCode: '',
    },
  };
  return (
    <Modal>
      <div>
        <Tag>HHG</Tag>
        <ModalClose handleClick={() => closeModal()} />
      </div>
      <ModalTitle>{title}</ModalTitle>
      <Formik onSubmit={(e) => onSave(e)} initialValues={initialValues}>
        {({ isValid }) => {
          return (
            <Form>
              <Label htmlFor="officeRemarks">Office remarks</Label>
              <Field as={Textarea} data-testid="officeRemarks" label="No" name="officeRemarks" id="officeRemarks" />
              <ModalActions>
                <Button type="submit">Save</Button>
                <Button type="button" secondary onClick={closeModal} disabled={!isValid}>
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

ServiceItemUpdateModal.propTypes = {
  closeModal: PropTypes.func.isRequired,
  onSave: PropTypes.func.isRequired,
  title: PropTypes.string.isRequired,
};

ServiceItemUpdateModal.displayName = 'ServiceItemUpdateModal';
export default connectModal(ServiceItemUpdateModal);
