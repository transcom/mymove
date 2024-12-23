import React from 'react';
import classnames from 'classnames';
import { PropTypes } from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Button } from '@trussworks/react-uswds';

import { AddressShape } from '../../../types/address';
import { formatAddress, formatCityStateAndPostalCode } from '../../../utils/shipmentDisplay';
import DataTableWrapper from '../../DataTableWrapper/index';
import DataTable from '../../DataTable/index';

import styles from './ShipmentAddresses.module.scss';

import { shipmentStatuses } from 'constants/shipments';
import { ShipmentOptionsOneOf, ShipmentStatusesOneOf } from 'types/shipment';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';

const ShipmentAddresses = ({
  pickupAddress,
  destinationAddress,
  originDutyLocation,
  destinationDutyLocation,
  handleShowDiversionModal,
  shipmentInfo,
  isMoveLocked,
  poeLocation,
}) => {
  let pickupHeader;
  let destinationHeader;
  switch (shipmentInfo.shipmentType) {
    case SHIPMENT_OPTIONS.HHG:
      pickupHeader = "Customer's addresses";
      destinationHeader = '';
      break;
    case SHIPMENT_OPTIONS.NTS:
      pickupHeader = 'Pickup Address';
      destinationHeader = 'Facility address';
      break;
    case SHIPMENT_OPTIONS.NTSR:
      pickupHeader = 'Facility address';
      destinationHeader = 'Delivery Address';
      break;
    default:
      pickupHeader = "Customer's addresses";
      destinationHeader = '';
  }

  return (
    <DataTableWrapper className={classnames('maxw-tablet', 'table--data-point-group', styles.mtoShipmentAddresses)}>
      <DataTable
        columnHeaders={[
          'Authorized addresses',
          <div className={styles.rightAlignButtonWrapper}>
            {shipmentInfo.status !== shipmentStatuses.CANCELED &&
              shipmentInfo.status !== shipmentStatuses.CANCELLATION_REQUESTED &&
              shipmentInfo.shipmentType !== SHIPMENT_OPTIONS.PPM && (
                <Restricted to={permissionTypes.createShipmentDiversionRequest}>
                  <Restricted to={permissionTypes.updateMTOPage}>
                    <Button
                      type="button"
                      onClick={() => handleShowDiversionModal(shipmentInfo)}
                      unstyled
                      disabled={isMoveLocked}
                    >
                      Request Diversion
                    </Button>
                  </Restricted>
                </Restricted>
              )}
          </div>,
        ]}
        dataRow={[
          formatCityStateAndPostalCode(originDutyLocation),
          formatCityStateAndPostalCode(destinationDutyLocation),
        ]}
        icon={<FontAwesomeIcon icon="arrow-right" />}
      />
      <DataTable
        columnHeaders={[pickupHeader, destinationHeader]}
        dataRow={[
          pickupAddress ? formatAddress(pickupAddress) : '—',
          destinationAddress ? formatAddress(destinationAddress) : '—',
        ]}
        icon={<FontAwesomeIcon icon="arrow-right" />}
        data-testid="pickupDestinationAddress"
      />
      <DataTable columnHeaders={['Port Location']} dataRow={[poeLocation]} />
    </DataTableWrapper>
  );
};

ShipmentAddresses.propTypes = {
  pickupAddress: AddressShape,
  destinationAddress: AddressShape,
  originDutyLocation: AddressShape,
  destinationDutyLocation: AddressShape,
  handleShowDiversionModal: PropTypes.func.isRequired,
  shipmentInfo: PropTypes.shape({
    id: PropTypes.string.isRequired,
    eTag: PropTypes.string.isRequired,
    status: ShipmentStatusesOneOf.isRequired,
    shipmentType: ShipmentOptionsOneOf.isRequired,
  }).isRequired,
};

ShipmentAddresses.defaultProps = {
  pickupAddress: {},
  destinationAddress: {},
  originDutyLocation: {},
  destinationDutyLocation: {},
};

export default ShipmentAddresses;
