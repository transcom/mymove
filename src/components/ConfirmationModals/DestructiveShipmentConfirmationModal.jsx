import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const DestructiveShipmentConfirmationModal = ({
  onClose,
  onSubmit,
  shipmentID,
  title,
  content,
  submitText,
  closeText,
}) => (
  <Modal>
    <ModalClose handleClick={() => onClose()} />
    <ModalTitle>
      <h3>{title}</h3>
    </ModalTitle>
    <p>{content}</p>
    <ModalActions autofocus="true">
      <Button data-focus="true" className="usa-button--destructive" type="submit" onClick={() => onSubmit(shipmentID)}>
        {submitText}
      </Button>
      <Button className="usa-button--secondary" type="button" onClick={() => onClose()} data-testid="modalBackButton">
        {closeText}
      </Button>
    </ModalActions>
  </Modal>
);

DestructiveShipmentConfirmationModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,

  shipmentID: PropTypes.string.isRequired,

  title: PropTypes.string,
  content: PropTypes.string,
  submitText: PropTypes.string,
  closeText: PropTypes.string,
};

DestructiveShipmentConfirmationModal.defaultProps = {
  title: 'Are you sure?',
  content:
    'You’ll lose all the information in this shipment. If you want it back later, you’ll have to request a new shipment.',
  submitText: 'Delete shipment',
  closeText: 'Keep shipment',
};

DestructiveShipmentConfirmationModal.displayName = 'DestructiveShipmentConfirmationModal';

export default connectModal(DestructiveShipmentConfirmationModal);
