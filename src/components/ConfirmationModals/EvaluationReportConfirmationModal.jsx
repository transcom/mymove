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
*/

export const EvaluationReportConfirmationModal = ({
  modalTitle,
  modalClose,
  reportId,
  moveCode,
  customerInfo,
  grade,
  shipmentId,
  closeModalOptions,
  submitModalOptions,
}) => (
  <Modal>
    {modalClose && <ModalClose handleClick={closeModal} />}
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
        <Button data-focus="true" className="usa-button--destructive" type="submit" onClick={closeModalOptions.onClick}>
          {closeModalOptions.text}
        </Button>
      )}
      {submitModalOptions && (
        <Button
          className="usa-button--secondary"
          type="button"
          onClick={submitModalOptions.onClick}
          data-testid="modalBackButton"
        >
          {submitModalOptions.text}
        </Button>
      )}
    </ModalActions>
  </Modal>
);

EvaluationReportConfirmationModal.PropTypes = {
  modalTitle: PropTypes.html,
  reportId: PropTypes.string.isRequired,
  moveCode: PropTypes.string.isRequired,
  customerInfo: CustomerShape.isRequired,
  grade: PropTypes.string.isRequired,
  shipmentId: PropTypes.string,
  closeModalOptions: PropTypes.shape({
    onClick: PropTypes.func,
    text: PropTypes.string,
  }),
  submitModalOptions: PropTypes.shape({
    onclick: PropTypes.func,
    text: PropTypes.string,
  }),
};

EvaluationReportConfirmationModal.defaultProps = {
  modalTitle: null,
  shipmentId: '',
  closeModalOptions: null,
  submitModalOptions: null,
};

EvaluationReportConfirmationModal.displayName = 'EvaluationReportConfirmationModal';

export default connectModal(EvaluationReportConfirmationModal);
