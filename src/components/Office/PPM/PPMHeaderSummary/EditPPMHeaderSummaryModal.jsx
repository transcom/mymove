import React from 'react';
import { Formik } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button } from '@trussworks/react-uswds';

import styles from './EditPPMHeaderSummaryModal.module.scss';

import { formatCentsTruncateWhole } from 'utils/formatters';
import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import { DatePickerInput } from 'components/form/fields';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';

const EditPPMHeaderSummaryModal = ({ sectionType, sectionInfo, onClose, onSubmit, editSectionName }) => {
  const { actualMoveDate, advanceAmountReceived } = sectionInfo;
  let title = 'Edit';
  if (sectionType === 'shipmentInfo') {
    title = 'Edit Shipment Info';
  } else if (sectionType === 'incentives') {
    title = 'Edit Incentives/Costs';
  }
  const initialValues = {
    editSectionName,
    actualMoveDate: actualMoveDate || '',
    advanceAmountReceived: formatCentsTruncateWhole(advanceAmountReceived).replace(/,/g, ''),
  };

  const validationSchema = Yup.object().shape({
    actualMoveDate: Yup.date()
      .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
      .when('editSectionName', {
        is: 'actualMoveDate',
        then: (schema) => schema.required('Required').max(new Date(), 'Date cannot be in the future'),
      }),
    advanceAmountReceived: Yup.number().when('editSectionName', {
      is: 'advanceAmountReceived',
      then: (schema) => schema.required('Required'),
    }),
  });

  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.EditPPMHeaderSummaryModal}>
          <ModalClose handleClick={() => onClose()} />
          <ModalTitle className={styles.ModalTitle}>
            <h3>{title}</h3>
          </ModalTitle>
          <Formik validationSchema={validationSchema} initialValues={initialValues} onSubmit={onSubmit}>
            {({ isValid }) => {
              return (
                <Form>
                  <div>
                    {editSectionName === 'actualMoveDate' && (
                      <DatePickerInput
                        name="actualMoveDate"
                        label="Actual move start date"
                        id="actualMoveDate"
                        disabledDays={{ after: new Date() }}
                      />
                    )}
                    {editSectionName === 'advanceAmountReceived' && (
                      <MaskedTextField
                        label="Advance received"
                        name="advanceAmountReceived"
                        id="advanceAmountReceived"
                        defaultValue="0"
                        mask={Number}
                        scale={0} // digits after point, 0 for integers
                        signed={false} // disallow negative
                        thousandsSeparator=","
                        lazy={false} // immediate masking evaluation
                        prefix="$"
                      />
                    )}
                  </div>
                  <ModalActions>
                    <Button type="submit" disabled={!isValid}>
                      Save
                    </Button>
                    <Button
                      type="button"
                      onClick={() => onClose()}
                      data-testid="modalCancelButton"
                      outline
                      className={styles.CancelButton}
                    >
                      Cancel
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
};

EditPPMHeaderSummaryModal.propTypes = {
  sectionType: PropTypes.string.isRequired,
  sectionInfo: PropTypes.shape({
    actualMoveDate: PropTypes.string,
    advanceAmountReceived: PropTypes.number,
  }),
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  editSectionName: PropTypes.string.isRequired,
};

EditPPMHeaderSummaryModal.defaultProps = {
  sectionInfo: {
    actualMoveDate: '',
    advanceAmountReceived: 0,
  },
};
export default EditPPMHeaderSummaryModal;
