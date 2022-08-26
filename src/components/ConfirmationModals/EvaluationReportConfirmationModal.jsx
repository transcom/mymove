import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';
import { CustomerShape } from 'types';

import EvaluationReportContainer from 'components/Office/EvaluationReportTable/EvaluationReportContainer';

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
  reportId,
  moveCode,
  customerInfo,
  grade,
  shipmentId,
  closeModalOptions,
  submitModalOptions,
}) => (
  <Modal>
    {modalTopRightClose && <ModalClose handleClick={modalTopRightClose} dataTestId={'modalCloseButtonTop'} />}
    {modalTitle && <ModalTitle>{modalTitle}</ModalTitle>}
    <EvaluationReportContainer
      evaluationReportId={reportId}
      moveCode={moveCode}
      customerInfo={customerInfo}
      grade={grade}
      shipmentId={shipmentId}
    />
    <ModalActions autofocus="true">
      {closeModalOptions && (
        <ModalClose>
          handleClick={closeModalOptions.handleClick}
          buttonContent={closeModalOptions.content}
          dataTestId={'modalCloseButtonBottom'}
        </ModalClose>
      )}
      {submitModalOptions && (
        <ModalSubmit>
          handleClick={submitModalOptions.handleClick}
          buttonContent={submitModalOptions.content}
        </ModalSubmit>
      )}
    </ModalActions>
  </Modal>
);

// TODO: check that shipmentId, reportId types make sense
// Why is there this typo error here?
EvaluationReportConfirmationModal.propTypes = {
  modalTitle: PropTypes.element,
  modalTopRightClose: PropTypes.func,
  reportId: PropTypes.string.isRequired,
  moveCode: PropTypes.string.isRequired,
  customerInfo: CustomerShape.isRequired,
  grade: PropTypes.string.isRequired,
  shipmentId: PropTypes.string,
  closeModalOptions: PropTypes.shape({
    handleClick: PropTypes.func,
    content: PropTypes.string,
  }),
  submitModalOptions: PropTypes.shape({
    handleClick: PropTypes.func,
    text: PropTypes.string,
  }),
};

EvaluationReportConfirmationModal.defaultProps = {
  modalTitle: <></>,
  modalTopRightClose: () => {},
  shipmentId: '',
  closeModalOptions: null,
  submitModalOptions: null,
};

EvaluationReportConfirmationModal.displayName = 'EvaluationReportConfirmationModal';

export default connectModal(EvaluationReportConfirmationModal);
