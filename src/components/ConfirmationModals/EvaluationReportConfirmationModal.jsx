import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './EvaluationReportConfirmationModal.module.scss';

import Modal, { ModalTitle, ModalClose, connectModal } from 'components/Modal/Modal';
import { CustomerShape, EvaluationReportShape, ShipmentShape } from 'types';
import EvaluationReportPreview from 'components/Office/EvaluationReportPreview/EvaluationReportPreview';

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
    *

  modalActions - optional
    * Element to render at the bottom of the modal, typically one or more cancle/submit type buttons.
    * Note: There are several unique button layout/styles used. Passed as a prop to allow for flexibility.
*/

export const EvaluationReportConfirmationModal = ({
  modalTitle,
  modalTopRightClose,
  modalActions,
  evaluationReport,
  moveCode,
  customerInfo,
  grade,
  mtoShipments,
  className,
  bordered,
}) => (
  <Modal className={classnames(styles.evaluationReportModal, className)}>
    {modalTopRightClose && <ModalClose handleClick={() => modalTopRightClose()} data-testid="modalCloseButtonTop" />}
    {modalTitle && <ModalTitle className={styles.titleSection}>{modalTitle}</ModalTitle>}
    <EvaluationReportPreview
      moveCode={moveCode}
      mtoShipments={mtoShipments}
      customerInfo={customerInfo}
      grade={grade}
      evaluationReport={evaluationReport}
      bordered={bordered}
    />
    {modalActions}
  </Modal>
);

EvaluationReportConfirmationModal.propTypes = {
  modalTitle: PropTypes.element,
  modalTopRightClose: PropTypes.func,
  evaluationReport: EvaluationReportShape.isRequired,
  mtoShipments: PropTypes.arrayOf(ShipmentShape),
  moveCode: PropTypes.string.isRequired,
  customerInfo: CustomerShape.isRequired,
  grade: PropTypes.string.isRequired,
  className: PropTypes.string,
  bordered: PropTypes.bool,
  modalActions: PropTypes.element,
};

EvaluationReportConfirmationModal.defaultProps = {
  modalTitle: null,
  modalTopRightClose: null,
  mtoShipments: null,
  className: null,
  bordered: false,
  modalActions: null,
};

EvaluationReportConfirmationModal.displayName = 'EvaluationReportConfirmationModal';

export default connectModal(EvaluationReportConfirmationModal);
