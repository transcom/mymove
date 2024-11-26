import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from '../BoatShipmentForm/BoatShipmentForm.module.scss';

import { boatShipmentTypes } from 'constants/shipments';
import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from 'components/Modal/Modal';

const boatConfirmationMessage = (isDimensionsMeetReq, boatShipmentType, isEditPage) => {
  let header = 'Boat Shipment';
  let message = '';

  // does not meet requirement to be a boat shipment
  if (!isDimensionsMeetReq) {
    // delete boat shipment and move to HHG
    if (isEditPage) {
      header = 'Movers pack and ship it, paid by the government (HHG)';
      message = (
        <p>
          Your boat meets the requirements to be moved with your HHG shipment. Click &quot;Delete & Continue&quot; to
          remove the Boat shipment and complete your shipment as HHG.
        </p>
      );
    } else {
      // move to HHG
      header = 'Movers pack and ship it, paid by the government (HHG)';
      message = (
        <p>
          Your boat meets the requirements to be moved with your HHG shipment. Click &quot;Continue&quot; to complete
          your shipment as HHG.
        </p>
      );
    }
  } else {
    // meets the requirement to be a boat shipment
    switch (boatShipmentType) {
      case boatShipmentTypes.TOW_AWAY:
        header = 'Boat Tow-Away (BTA)';
        message = (
          <p>
            Your boat qualifies to move as its own shipment and has an accompanying trailer that can be used to tow it
            to your delivery address, a Boat Tow-Away (BTA) shipment. Click &quot;Continue&quot; to proceed.
          </p>
        );
        break;
      case boatShipmentTypes.HAUL_AWAY:
        header = 'Boat Haul-Away (BHA)';
        message = (
          <p>
            Your boat qualifies to move as its own shipment and requires additional equipment to haul it to your
            delivery address, a Boat Haul-Away (BHA) shipment. Click &quot;Continue&quot; to proceed.
          </p>
        );
        break;
      default:
        break;
    }
  }

  return { header, message };
};

export const BoatShipmentConfirmationModal = ({
  isDimensionsMeetReq,
  boatShipmentType,
  closeModal,
  handleConfirmationContinue,
  handleConfirmationRedirect,
  handleConfirmationDeleteAndRedirect,
  isSubmitting,
  isEditPage,
}) => {
  const { header, message } = boatConfirmationMessage(isDimensionsMeetReq, boatShipmentType, isEditPage);

  return (
    <Modal>
      <ModalClose handleClick={closeModal} />
      <ModalTitle>
        <h3>{header}</h3>
      </ModalTitle>
      {message}
      <ModalActions>
        <div className={styles.buttonContainer}>
          <Button
            data-testid="boatConfirmationBack"
            className={styles.backButton}
            type="button"
            onClick={closeModal}
            secondary
            outline
          >
            Back
          </Button>
          {isEditPage && !isDimensionsMeetReq ? (
            <Button
              data-testid="boatConfirmationContinue"
              className="usa-button--destructive"
              type="button"
              onClick={handleConfirmationDeleteAndRedirect}
              disabled={isSubmitting}
            >
              Delete & Continue
            </Button>
          ) : (
            <Button
              data-testid="boatConfirmationContinue"
              className={styles.saveButton}
              type="button"
              onClick={isDimensionsMeetReq ? handleConfirmationContinue : handleConfirmationRedirect}
              disabled={isSubmitting}
            >
              Continue
            </Button>
          )}
        </div>
      </ModalActions>
    </Modal>
  );
};

BoatShipmentConfirmationModal.propTypes = {
  isDimensionsMeetReq: PropTypes.bool.isRequired,
  boatShipmentType: PropTypes.string.isRequired,
  closeModal: PropTypes.func,
  handleConfirmationContinue: PropTypes.func.isRequired,
  handleConfirmationRedirect: PropTypes.func.isRequired,
  handleConfirmationDeleteAndRedirect: PropTypes.func.isRequired,
  isSubmitting: PropTypes.bool,
};

BoatShipmentConfirmationModal.defaultProps = {
  closeModal: () => {},
  isSubmitting: false,
};

BoatShipmentConfirmationModal.displayName = 'BoatConfirmationModal';

export default connectModal(BoatShipmentConfirmationModal);
