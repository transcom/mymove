import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

export const MoveInfoModal = ({ closeModal, enablePPM, enableUB, hasOconusDutyLocation }) => (
  <Modal data-testid="moveInfoModal">
    <ModalClose handleClick={closeModal} />
    <ModalTitle>
      <h3 data-testid="moveInfoModalHeading">More info about shipments</h3>
    </ModalTitle>

    <h4>
      <strong>HHG: Professional movers pack and ship your things, the government pays</strong>
    </h4>
    <p data-testid="hhgSubHeading">The moving company works out details with you, but handles everything.</p>
    <h5>Pros</h5>
    <ul data-testid="hhgProsList">
      <li>Everything is packed and moved for you</li>
      <li>Expert movers care for your things</li>
      <li>Anything damaged in professional shipments will be replaced</li>
    </ul>
    <h5>Cons</h5>
    <ul data-testid="hhgConsList">
      <li>Can only move on weekdays</li>
      <li>May have to work around availability of movers</li>
    </ul>
    {enablePPM && (
      <>
        <h4>
          <strong>PPM: You get your things packed and moved, the government pays you</strong>
        </h4>
        <p data-testid="ppmSubHeading">You pack and move your own things, or arrange for someone else do it for you.</p>
        <h5>Pros</h5>
        <ul data-testid="ppmProsList">
          <li>Keep your things with you at all times</li>
          <li>Get paid for the weight you move</li>
          <li>Flexible dates, routes, timing</li>
          <li>You can hire movers, equipment, or portable storage</li>
        </ul>
        <h5>Cons</h5>
        <ul data-testid="ppmConsList">
          <li>You pack and move everything</li>
          <li>You’re responsible if your things get damaged — no compensation</li>
          <li>The more you own, the more you have to do</li>
        </ul>
      </>
    )}
    {enableUB && hasOconusDutyLocation && (
      <>
        <h4>
          <strong>
            UB: Professional movers pack and ship your more essential personal property, the government pays
          </strong>
        </h4>
        <p data-testid="ubSubHeading">The moving company works out details with you, but handles everything.</p>
        <h5>Pros</h5>
        <ul data-testid="ubProsList">
          <li>Everything is packed and moved for you</li>
          <li>Expert movers care for your things</li>
          <li>Anything damaged in professional shipments will be replaced</li>
          <li>Essential items are packed as a separate shipment</li>
          <li>Shorter allowable transit time than a standard HHG shipment; should arrive sooner at your destination</li>
        </ul>
        <h5>Cons</h5>
        <ul data-testid="ubConsList">
          <li>Can only move on weekdays</li>
          <li>May have to work around availability of movers</li>
          <li>Your UB shipment has its own weight limitation</li>
          <li>Only certain kinds of personal property are allowed in a UB shipment (check with your counselor)</li>
        </ul>
      </>
    )}
    <ModalActions>
      <Button secondary type="button" onClick={closeModal}>
        Got it
      </Button>
    </ModalActions>
  </Modal>
);

MoveInfoModal.propTypes = {
  closeModal: PropTypes.func,
  enablePPM: PropTypes.bool,
  enableUB: PropTypes.bool,
  hasOconusDutyLocation: PropTypes.bool,
};

MoveInfoModal.defaultProps = {
  closeModal: () => {},
  enablePPM: true,
  enableUB: true,
  hasOconusDutyLocation: true,
};

MoveInfoModal.displayName = 'MoveInfoModal';

export default connectModal(MoveInfoModal);
