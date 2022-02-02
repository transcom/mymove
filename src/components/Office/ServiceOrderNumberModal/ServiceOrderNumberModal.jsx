import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import Modal, { connectModal, ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
import { Form } from 'components/form';
import TextField from 'components/form/fields/TextField/TextField';

const validationSchema = Yup.object().shape({
  serviceOrderNumber: Yup.string()
    .matches(/^[0-9a-zA-Z]+$/, 'Letters and numbers only')
    .required('Required'),
});

const ServiceOrderNumberModal = ({ onClose, onSubmit, serviceOrderNumber }) => {
  const handleFormSubmit = (values) => onSubmit(values);

  return (
    <div data-testid="ServiceOrderNumber">
      <Modal>
        <ModalClose handleClick={onClose} />

        <ModalTitle>
          <h2>Edit service order number</h2>
        </ModalTitle>

        <Formik initialValues={{ serviceOrderNumber }} onSubmit={handleFormSubmit} validationSchema={validationSchema}>
          <Form>
            <TextField label="Service order number" id="facilityServiceOrderNumber" name="serviceOrderNumber" />

            <ModalActions>
              <Button type="submit">Save</Button>
              <Button type="button" secondary onClick={onClose}>
                Cancel
              </Button>
            </ModalActions>
          </Form>
        </Formik>
      </Modal>
    </div>
  );
};

ServiceOrderNumberModal.propTypes = {
  onClose: PropTypes.func,
  onSubmit: PropTypes.func,
  serviceOrderNumber: PropTypes.string,
};

ServiceOrderNumberModal.defaultProps = {
  onClose: () => {},
  onSubmit: () => {},
  serviceOrderNumber: '',
};

export default connectModal(ServiceOrderNumberModal);
