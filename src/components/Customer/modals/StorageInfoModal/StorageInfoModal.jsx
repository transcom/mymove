import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const StorageInfoModal = ({ closeModal }) => (
  <Modal>
    <ModalClose handleClick={closeModal} />
    <ModalTitle>
      <h3>Long-term storage info</h3>
    </ModalTitle>
    <h4>
      <strong>Long-term storage (NTS)</strong>
    </h4>
    <p>
      Put some or all of your things into storage as part of one move, and get it out of storage on a future move. Your
      move counselor can verify whether or not you qualify to put things into long-term storage on this move.
    </p>
    <ul>
      <li>The weight of this shipment counts against your weight allowance</li>
      <li>Useful when you can’t take all your things to your new location</li>
      <li>Common in OCONUS moves, but may not be available in CONUS</li>
      <li>Stored in a government-approved facility, typically near your starting location</li>
    </ul>
    <p>
      NTS (short for “non-temp storage”) lasts 6 months or longer. Do not count on easy access to things in storage. You
      can retrieve them during a future move.
    </p>

    <h4>
      <strong>Long-term storage release (NTS-R)</strong>
    </h4>
    <p>
      Retrieval of things you placed into long-term storage (NTS) during an earlier move. If you don’t have things in
      storage, you don’t need to request an NTS-R.
    </p>
    <p>Movers pick up your things at the storage facility and deliver them to your new location.</p>
    <ul>
      <li>The weight of this shipment counts against your weight allowance</li>
      <li>It may not be possible to do a partial release from storage</li>
      <li>Will arrive as its own shipment, not as part of any other shipments you select</li>
      <li>If you don’t know where your shipment was stored, talk to your move counselor — they can look it up</li>
    </ul>
    <p>This is known as an NTS-R — a “non-temporary storage release.”</p>
    <ModalActions>
      <Button secondary type="button" onClick={closeModal}>
        Got it
      </Button>
    </ModalActions>
  </Modal>
);

StorageInfoModal.propTypes = {
  closeModal: PropTypes.func,
};

StorageInfoModal.defaultProps = {
  closeModal: () => {},
};

StorageInfoModal.displayName = 'StorageInfoModal';

export default connectModal(StorageInfoModal);
