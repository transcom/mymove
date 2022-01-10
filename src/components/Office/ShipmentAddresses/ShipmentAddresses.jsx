import React from 'react';
import classnames from 'classnames';
import { PropTypes } from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Button } from '@trussworks/react-uswds';

import { AddressShape } from '../../../types/address';
import { formatAddress } from '../../../utils/shipmentDisplay';
import DataTableWrapper from '../../DataTableWrapper/index';
import DataTable from '../../DataTable/index';

import styles from './ShipmentAddresses.module.scss';

import { shipmentStatuses } from 'constants/shipments';
import { ShipmentOptionsOneOf, ShipmentStatusesOneOf } from 'types/shipment';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const ShipmentAddresses = ({
  pickupAddress,
  destinationAddress,
  originDutyStation,
  destinationDutyStation,
  handleDivertShipment,
  shipmentInfo,
}) => {
  let pickupHeader;
  let destinationHeader;
  switch (shipmentInfo.shipmentType) {
    case SHIPMENT_OPTIONS.HHG:
      pickupHeader = "Customer's addresses";
      destinationHeader = '';
      break;
    case SHIPMENT_OPTIONS.NTS:
      pickupHeader = 'Pickup address';
      destinationHeader = 'Facility address';
      break;
    case SHIPMENT_OPTIONS.NTSR:
      pickupHeader = 'Facility address';
      destinationHeader = 'Delivery address';
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
            {shipmentInfo.shipmentStatus !== shipmentStatuses.CANCELED && (
              <Button
                type="button"
                onClick={() => handleDivertShipment(shipmentInfo.shipmentID, shipmentInfo.ifMatchEtag)}
                unstyled
              >
                Request diversion
              </Button>
            )}
          </div>,
        ]}
        dataRow={[formatAddress(originDutyStation), formatAddress(destinationDutyStation)]}
        icon={<FontAwesomeIcon icon="arrow-right" />}
      />
      <DataTable
        columnHeaders={[pickupHeader, destinationHeader]}
        dataRow={[
          pickupAddress ? formatAddress(pickupAddress) : '—',
          destinationAddress ? formatAddress(destinationAddress) : '—',
        ]}
        icon={<FontAwesomeIcon icon="arrow-right" />}
      />
    </DataTableWrapper>
  );
};

ShipmentAddresses.propTypes = {
  pickupAddress: AddressShape,
  destinationAddress: AddressShape,
  originDutyStation: AddressShape,
  destinationDutyStation: AddressShape,
  handleDivertShipment: PropTypes.func.isRequired,
  shipmentInfo: PropTypes.shape({
    shipmentID: PropTypes.string.isRequired,
    ifMatchEtag: PropTypes.string.isRequired,
    shipmentStatus: ShipmentStatusesOneOf.isRequired,
    shipmentType: ShipmentOptionsOneOf.isRequired,
  }).isRequired,
};

ShipmentAddresses.defaultProps = {
  pickupAddress: {},
  destinationAddress: {},
  originDutyStation: {},
  destinationDutyStation: {},
};

export default ShipmentAddresses;
