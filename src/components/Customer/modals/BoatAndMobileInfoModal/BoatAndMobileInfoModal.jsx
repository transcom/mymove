import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const BoatAndMobileInfoModal = ({ closeModal }) => (
  <Modal>
    <ModalClose handleClick={closeModal} />
    <ModalTitle>
      <h3>Boat & mobile homes info</h3>
    </ModalTitle>
    <h4>
      <strong>Boat shipment</strong>
    </h4>
    <p>
      This option is for privately owned boats or personal watercraft (including but not limited to canoes, kayaks,
      dinghies, row boats, jet skis, and sculls) with or without an associated trailer; over 14 feet in length or over 6
      feet 10 inches in width or over 6 feet 5 inches in height. Boats or personal watercraft equal to or less than the
      above dimensions shall be shipped with household goods and not be considered a separate shipment. If your boat is
      under those dimensions, choose the &quot;HHG&quot; option above.
    </p>
    <h4>
      <strong>Mobile Home shipment</strong>
    </h4>
    <p>This option is for privately owned mobile homes.</p>
    <ModalActions>
      <Button secondary type="button" onClick={closeModal}>
        Got it
      </Button>
    </ModalActions>
  </Modal>
);

BoatAndMobileInfoModal.propTypes = {
  closeModal: PropTypes.func,
};

BoatAndMobileInfoModal.defaultProps = {
  closeModal: () => {},
};

BoatAndMobileInfoModal.displayName = 'BoatAndMobileInfoModal';

export default connectModal(BoatAndMobileInfoModal);
