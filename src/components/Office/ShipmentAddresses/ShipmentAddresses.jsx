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
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';

const ShipmentAddresses = ({
  pickupAddress,
  destinationAddress,
  originDutyLocation,
  destinationDutyLocation,
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
            {shipmentInfo.status !== shipmentStatuses.CANCELED && (
              <Restricted to={permissionTypes.createShipmentDiversionRequest}>
                <Button type="button" onClick={() => handleDivertShipment(shipmentInfo.id, shipmentInfo.eTag)} unstyled>
                  Request diversion
                </Button>
              </Restricted>
            )}
          </div>,
        ]}
        dataRow={[formatAddress(originDutyLocation), formatAddress(destinationDutyLocation)]}
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
    </DataTableWrapper>
  );
};

ShipmentAddresses.propTypes = {
  pickupAddress: AddressShape,
  destinationAddress: AddressShape,
  originDutyLocation: AddressShape,
  destinationDutyLocation: AddressShape,
  handleDivertShipment: PropTypes.func.isRequired,
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
