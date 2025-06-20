import React from 'react';
import PropTypes from 'prop-types';
import { Button, Form } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import * as Yup from 'yup';

import styles from './TerminateShipmentModal.module.scss';

import Modal, { ModalTitle, ModalActions, connectModal, ModalClose } from 'components/Modal/Modal';
import TextField from 'components/form/fields/TextField/TextField';
import Hint from 'components/Hint';

export const TerminateShipmentModal = ({ onClose, onSubmit, shipmentID, shipmentLocator }) => {
  const validationSchema = Yup.object().shape({
    terminationComments: Yup.string().required('Required'),
  });

  const initialValues = {
    terminationComments: '',
  };

  return (
    <Modal className={styles.modal}>
      <ModalClose handleClick={() => onClose()} />
      <ModalTitle className={styles.center}>
        <h3>Shipment termination</h3>
        <Hint>{shipmentLocator}</Hint>
      </ModalTitle>
      <div>
        <Formik initialValues={{ ...initialValues }} validationSchema={validationSchema} validateOnMount>
          {({ values, isValid, isSubmitting }) => {
            return (
              <Form>
                <TextField
                  data-testid="terminationComments"
                  label="Termination reason"
                  id="terminationComments"
                  name="terminationComments"
                  prefix="TERMINATED FOR CAUSE:"
                  required
                  showRequiredAsterisk
                />
                <ModalActions>
                  <Button
                    className="usa-button--secondary"
                    type="button"
                    data-testid="modalBackBtn"
                    onClick={() => onClose()}
                  >
                    Cancel
                  </Button>
                  <Button
                    type="submit"
                    data-focus="true"
                    data-testid="modalSubmitBtn"
                    disabled={!isValid || isSubmitting}
                    onClick={() => onSubmit(shipmentID, values)}
                  >
                    Terminate
                  </Button>
                </ModalActions>
              </Form>
            );
          }}
        </Formik>
      </div>
    </Modal>
  );
};

TerminateShipmentModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  shipmentID: PropTypes.string.isRequired,
  shipmentLocator: PropTypes.string.isRequired,
};

TerminateShipmentModal.displayName = 'TerminateShipmentModal';

export default connectModal(TerminateShipmentModal);
