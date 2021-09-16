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
import { ShipmentStatusesOneOf } from 'types/shipment';

const ShipmentAddresses = ({
  pickupAddress,
  destinationAddress,
  originDutyStation,
  destinationDutyStation,
  handleDivertShipment,
  shipmentInfo,
}) => {
  return (
    <DataTableWrapper className={classnames('maxw-tablet', styles.mtoShipmentAddresses)}>
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
        columnHeaders={["Customer's addresses", '']}
        dataRow={[formatAddress(pickupAddress), formatAddress(destinationAddress)]}
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
  }).isRequired,
};

ShipmentAddresses.defaultProps = {
  pickupAddress: {},
  destinationAddress: {},
  originDutyStation: {},
  destinationDutyStation: {},
};

export default ShipmentAddresses;
