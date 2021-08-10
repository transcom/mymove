import React from 'react';
import { Formik } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button } from '@trussworks/react-uswds';

import styles from './EditMaxBillableWeightModal.module.scss';

import MaskedTextField from 'components/form/fields/MaskedTextField';
import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
import { formatWeight } from 'shared/formatters';

const maxBillableWeightSchema = Yup.object().shape({
  maxBillableWeight: Yup.number().required('Required'),
});

const EditMaxBillableWeightModal = ({ onClose, onSubmit, defaultWeight, maxBillableWeight }) => (
  <div className={styles.EditMaxBillableWeightModal}>
    <Overlay />
    <ModalContainer>
      <Modal>
        <ModalClose className={styles.weightModalClose} handleClick={() => onClose()} />
        <ModalTitle>
          <h4>Edit max billable weight</h4>
        </ModalTitle>
        <dl>
          <dt>Default: </dt>
          <dd>{formatWeight(defaultWeight)}</dd>
        </dl>
        <Formik initialValues={{ maxBillableWeight }} validationSchema={maxBillableWeightSchema} onSubmit={onSubmit}>
          {({ isValid }) => {
            return (
              <Form>
                <MaskedTextField name="maxBillableWeight" label="New max billable weight" />
                <ModalActions>
                  <Button type="submit" disabled={!isValid}>
                    Save
                  </Button>
                  <Button
                    type="button"
                    onClick={() => onClose()}
                    className={styles.backButton}
                    data-testid="modalBackButton"
                    outline
                  >
                    Back
                  </Button>
                </ModalActions>
              </Form>
            );
          }}
        </Formik>
      </Modal>
    </ModalContainer>
  </div>
);

EditMaxBillableWeightModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  defaultWeight: PropTypes.number.isRequired,
  maxBillableWeight: PropTypes.number,
};

EditMaxBillableWeightModal.defaultProps = {
  maxBillableWeight: undefined,
};

export default EditMaxBillableWeightModal;
