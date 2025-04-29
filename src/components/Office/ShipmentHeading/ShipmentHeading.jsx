import React, { useEffect, useState } from 'react';
import classNames from 'classnames';
import { PropTypes } from 'prop-types';
import { Button, Tag } from '@trussworks/react-uswds';
import { connect } from 'react-redux';
import { useQueryClient, useMutation } from '@tanstack/react-query';

import { AddressShape } from '../../../types/address';
import TerminateShipmentModal from '../TerminateShipmentModal/TerminateShipmentModal';

import styles from './shipmentHeading.module.scss';

import { shipmentStatuses } from 'constants/shipments';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';
import { roleTypes } from 'constants/userRoles';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { FEATURE_FLAG_KEYS } from 'shared/constants';
import { terminateShipment } from 'services/ghcApi';
import { MTO_SHIPMENTS } from 'constants/queryKeys';
import { milmoveLogger } from 'utils/milmoveLog';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';

function ShipmentHeading({ shipmentInfo, handleShowCancellationModal, isMoveLocked, activeRole, setFlashMessage }) {
  const [terminatingShipmentsFF, setTerminatingShipmentsFF] = useState(false);
  const [isShipmentTerminationModalVisible, setIsShipmentTerminationModalVisible] = useState(false);
  const { shipmentStatus } = shipmentInfo;
  // cancelation modal is visible if shipment is not already canceled, AND if shipment cancellation hasn't already been requested
  const showRequestCancellation =
    shipmentStatus !== shipmentStatuses.CANCELED && shipmentStatus !== shipmentStatuses.CANCELLATION_REQUESTED;
  const isCancellationRequested = shipmentStatus === shipmentStatuses.CANCELLATION_REQUESTED;
  const isDisabled = isMoveLocked || shipmentStatus === shipmentStatuses.TERMINATED_FOR_CAUSE;

  const queryClient = useQueryClient();
  const { mutate: mutateShipmentTermination } = useMutation(terminateShipment, {
    onSuccess: (updatedMTOShipment) => {
      setFlashMessage(
        `TERMINATION_SUCCESS_${updatedMTOShipment.id}`,
        'success',
        `Successfully terminated shipment ${updatedMTOShipment.shipmentLocator}`,
        '',
        true,
      );

      queryClient.invalidateQueries([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID]);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });

  useEffect(() => {
    const fetchData = async () => {
      setTerminatingShipmentsFF(await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.TERMINATING_SHIPMENTS));
    };
    fetchData();
  }, []);

  const canTerminate =
    !shipmentInfo.actualPickupDate &&
    shipmentInfo.shipmentStatus === shipmentStatuses.APPROVED &&
    terminatingShipmentsFF &&
    activeRole === roleTypes.CONTRACTING_OFFICER;

  const handleShipmentTerminationSubmit = (shipmentID, values) => {
    const body = {
      terminationReason: values.terminationComments,
    };
    mutateShipmentTermination({ shipmentID, body });

    setIsShipmentTerminationModalVisible(false);
  };

  const handleShipmentTerminationCancel = () => {
    setIsShipmentTerminationModalVisible(false);
  };

  return (
    <div className={classNames(styles.shipmentHeading, 'shipment-heading')}>
      <TerminateShipmentModal
        isOpen={isShipmentTerminationModalVisible}
        onClose={handleShipmentTerminationCancel}
        onSubmit={handleShipmentTerminationSubmit}
        shipmentID={shipmentInfo.shipmentID}
        shipmentLocator={shipmentInfo.shipmentLocator}
      />
      <div className={styles.shipmentHeadingType}>
        <h2>
          <span className={styles.marketCodeIndicator}>{shipmentInfo.marketCode}</span>
          {shipmentInfo.shipmentType}
        </h2>
        <div>
          {shipmentStatus === shipmentStatuses.TERMINATED_FOR_CAUSE && (
            <Tag className="usa-tag--cancellation">terminated for cause</Tag>
          )}
          {shipmentStatus === shipmentStatuses.CANCELED && <Tag className="usa-tag--cancellation">canceled</Tag>}
          {shipmentInfo.isDiversion && <Tag className="usa-tag--diversion">diversion</Tag>}
          {!shipmentInfo.isDiversion && shipmentStatus === shipmentStatuses.DIVERSION_REQUESTED && (
            <Tag className="usa-tag--diversion">diversion requested</Tag>
          )}
        </div>
      </div>
      <div>
        <h4>#{shipmentInfo.shipmentLocator}</h4>
      </div>
      <div className={styles.column}>
        {showRequestCancellation && (
          <Restricted to={permissionTypes.createShipmentCancellation}>
            <Restricted to={permissionTypes.updateMTOPage}>
              <Button
                data-testid="requestCancellationBtn"
                type="button"
                onClick={() => handleShowCancellationModal(shipmentInfo)}
                unstyled
                disabled={isDisabled}
              >
                Request Cancellation
              </Button>
            </Restricted>
          </Restricted>
        )}
        <Restricted to={permissionTypes.createShipmentTermination}>
          {canTerminate && (
            <Button
              data-testid="terminateShipmentBtn"
              type="button"
              onClick={() => {
                setIsShipmentTerminationModalVisible(true);
              }}
              unstyled
              disabled={isMoveLocked}
            >
              Terminate Shipment
            </Button>
          )}
        </Restricted>
        {isCancellationRequested && <Tag className="usa-tag--cancellation">Cancellation Requested</Tag>}
      </div>
    </div>
  );
}

ShipmentHeading.propTypes = {
  shipmentInfo: PropTypes.shape({
    shipmentID: PropTypes.string.isRequired,
    shipmentType: PropTypes.string.isRequired,
    isDiversion: PropTypes.bool,
    originCity: PropTypes.string.isRequired,
    originState: PropTypes.string.isRequired,
    originPostalCode: PropTypes.string.isRequired,
    destinationAddress: AddressShape,
    scheduledPickupDate: PropTypes.string.isRequired,
    shipmentStatus: PropTypes.string.isRequired,
    ifMatchEtag: PropTypes.string.isRequired,
    moveTaskOrderID: PropTypes.string.isRequired,
  }).isRequired,
  handleShowCancellationModal: PropTypes.func.isRequired,
};

const mapStateToProps = (state) => {
  return {
    activeRole: state.auth.activeRole,
  };
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(ShipmentHeading);
