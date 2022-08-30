import React from 'react';
import PropTypes from 'prop-types';

import styles from './EvaluationReportConfirmationModal.module.scss';

import Modal, { ModalTitle, ModalClose, ModalSubmit, ModalActions, connectModal } from 'components/Modal/Modal';
import { CustomerShape, EvaluationReportShape, ShipmentShape } from 'types';
import EvaluationReportContainer from 'components/Office/EvaluationReportPreview/EvaluationReportPreview';

/*
  Used for the Evaluation Report:
    Preview (and submit)
      Modal Title
      Close button on top right hand side
    View (already submitted report)
      No title/subtext
      Close button only on the bottom

  modalTitle - optional
    * format <html or react elements to render within the ModalTitle div>

  modalTopRightClose - optional
    * renders a close button in the top right corner
    * triggers the passed in onClick method

  reportId, moveCode - required
   * used to look up the report to fill in details

  grade - required
    ORDERS_RANK_OPTIONS[GRADE]

  shipmentId - required if a shipment eval report, optional for others

  close & submitModalOptions
    onClick
    content

*/

export const EvaluationReportConfirmationModal = ({
  modalTitle,
  modalTopRightClose,
  evaluationReport,
  moveCode,
  customerInfo,
  grade,
  mtoShipments,
  closeModalOptions,
  submitModalOptions,
}) => (
  <Modal className={styles.evaluationReportModal}>
    {modalTopRightClose && <ModalClose handleClick={() => modalTopRightClose()} data-testid="modalCloseButtonTop" />}
    {modalTitle && <ModalTitle className={styles.titleSection}>{modalTitle}</ModalTitle>}
    <EvaluationReportContainer
      evaluationReport={evaluationReport}
      moveCode={moveCode}
      mtoShipments={mtoShipments}
      customerInfo={customerInfo}
      grade={grade}
    />
    <ModalActions autofocus="true">
      {closeModalOptions && (
        <ModalClose
          handleClick={() => closeModalOptions.handleClick()}
          buttonContent={closeModalOptions.buttonContent}
          data-testid="modalCloseButtonBottom"
        />
      )}
      {submitModalOptions && (
        <ModalSubmit
          handleClick={() => submitModalOptions.handleClick()}
          buttonContent={submitModalOptions.buttonContent}
        />
      )}
    </ModalActions>
  </Modal>
);

// TODO: check that shipmentId, reportId types make sense
// Why is there this typo error here?
EvaluationReportConfirmationModal.propTypes = {
  modalTitle: PropTypes.element,
  modalTopRightClose: PropTypes.func,
  evaluationReport: EvaluationReportShape.isRequired,
  mtoShipments: PropTypes.arrayOf(ShipmentShape),
  moveCode: PropTypes.string.isRequired,
  customerInfo: CustomerShape.isRequired,
  grade: PropTypes.string.isRequired,
  closeModalOptions: PropTypes.shape({
    handleClick: PropTypes.func,
    buttonContent: PropTypes.string,
  }),
  submitModalOptions: PropTypes.shape({
    handleClick: PropTypes.func,
    buttonContent: PropTypes.string,
  }),
};

EvaluationReportConfirmationModal.defaultProps = {
  modalTitle: null,
  modalTopRightClose: null,
  mtoShipments: null,
  closeModalOptions: null,
  submitModalOptions: null,
};

EvaluationReportConfirmationModal.displayName = 'EvaluationReportConfirmationModal';

export default connectModal(EvaluationReportConfirmationModal);
