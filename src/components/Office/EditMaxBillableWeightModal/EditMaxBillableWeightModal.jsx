import React from 'react';
import { Formik } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button } from '@trussworks/react-uswds';

import styles from './EditMaxBillableWeightModal.module.scss';

import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { Form } from 'components/form';
import Modal, { ModalActions, ModalClose, ModalTitle, connectModal } from 'components/Modal/Modal';
import { formatWeight } from 'utils/formatters';
import { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

const maxBillableWeightSchema = Yup.object().shape({
  maxBillableWeight: Yup.number().min(1, 'Max billable weight must be greater than or equal to 1').required('Required'),
});

export const EditMaxBillableWeightModal = ({ onClose, onSubmit, defaultWeight, maxBillableWeight }) => {
  return (
    <Modal className={styles.EditMaxBillableWeightModal} onClose={() => onClose()}>
      <ModalClose className={styles.weightModalClose} handleClick={() => onClose()} />
      <ModalTitle>
        <h4>Edit max billable weight</h4>
      </ModalTitle>
      <dl>
        <dt>Default:</dt>
        <dd>{formatWeight(defaultWeight)}</dd>
      </dl>
      <Formik
        initialValues={{ maxBillableWeight: `${maxBillableWeight}` }}
        validationSchema={maxBillableWeightSchema}
        onSubmit={(values) => {
          onSubmit(Number.parseInt(values.maxBillableWeight, 10));
        }}
      >
        {({ isValid }) => {
          return (
            <Form>
              {requiredAsteriskMessage}
              <MaskedTextField
                name="maxBillableWeight"
                id="maxBillableWeight"
                label="New max billable weight"
                mask="num lbs"
                blocks={{
                  num: {
                    mask: Number,
                    signed: false,
                    scale: 0,
                    thousandsSeparator: ',',
                  },
                }}
                lazy={false}
                showRequiredAsterisk
                required
              />
              <ModalActions>
                <Button type="button" secondary onClick={() => onClose()} data-testid="modalBackButton" outline>
                  Back
                </Button>
                <Button type="submit" disabled={!isValid}>
                  Save
                </Button>
              </ModalActions>
            </Form>
          );
        }}
      </Formik>
    </Modal>
  );
};

EditMaxBillableWeightModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  defaultWeight: PropTypes.number.isRequired,
  maxBillableWeight: PropTypes.number,
};

EditMaxBillableWeightModal.defaultProps = {
  maxBillableWeight: undefined,
};

EditMaxBillableWeightModal.displayName = 'EditMaxBillableWeightModal';

export default connectModal(EditMaxBillableWeightModal);
