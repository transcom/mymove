import React from 'react';
import classnames from 'classnames';
import * as PropTypes from 'prop-types';

import styles from './importantShipmentDates.module.scss';

import DataTable from 'components/DataTable/index';
import DataTableWrapper from 'components/DataTableWrapper/index';
import { shipmentStatuses } from 'constants/shipments';
import { ShipmentOptionsOneOf, ShipmentStatusesOneOf } from 'types/shipment';

const ImportantShipmentDates = ({
  requestedPickupDate,
  plannedMoveDate,
  scheduledPickupDate,
  actualMoveDate,
  actualPickupDate,
  requiredDeliveryDate,
  requestedDeliveryDate,
  scheduledDeliveryDate,
  actualDeliveryDate,
  isPPM,
  shipmentInfo,
}) => {
  const headerPlannedMoveDate = isPPM ? 'Planned Move Date' : 'Requested pick up date';
  const headerActualMoveDate = isPPM ? 'Actual Move Date' : 'Scheduled pick up date';
  const headerActualPickupDate = isPPM ? '' : 'Actual pick up date';
  const emDash = '\u2014';
  return (
    <div className={classnames('maxw-tablet', styles.shipmentDatesContainer)}>
      <DataTableWrapper className="table--data-point-group">
        {shipmentInfo.isDiversion && (
          <DataTable columnHeaders={['Diversion Approved']} dataRow={[shipmentInfo.diversionReason || emDash]} />
        )}
        {!shipmentInfo.isDiversion && shipmentInfo.status === shipmentStatuses.DIVERSION_REQUESTED && (
          <DataTable columnHeaders={['Diversion Requested']} dataRow={[shipmentInfo.diversionReason || emDash]} />
        )}
        {!isPPM && <DataTable columnHeaders={['Required Delivery Date']} dataRow={[requiredDeliveryDate || emDash]} />}
        {!isPPM && (
          <DataTable
            columnHeaders={[headerPlannedMoveDate, headerActualMoveDate, headerActualPickupDate]}
            dataRow={[requestedPickupDate || emDash, scheduledPickupDate || emDash, actualPickupDate || emDash]}
          />
        )}
        {isPPM && (
          <DataTable
            columnHeaders={[headerPlannedMoveDate, headerActualMoveDate]}
            dataRow={[plannedMoveDate || emDash, actualMoveDate || emDash]}
          />
        )}
        {!isPPM && (
          <DataTable
            columnHeaders={['Requested delivery date', 'Scheduled delivery date', 'Actual delivery date']}
            dataRow={[requestedDeliveryDate || emDash, scheduledDeliveryDate || emDash, actualDeliveryDate || emDash]}
          />
        )}
      </DataTableWrapper>
    </div>
  );
};

ImportantShipmentDates.defaultProps = {
  requestedPickupDate: '',
  scheduledPickupDate: '',
  requiredDeliveryDate: '',
  actualPickupDate: '',
  plannedMoveDate: '',
  actualMoveDate: '',
  requestedDeliveryDate: '',
  scheduledDeliveryDate: '',
  actualDeliveryDate: '',
  isPPM: false,
  shipmentInfo: {
    id: '',
    status: '',
    shipmentType: '',
    isDiversion: false,
    diversionReason: '',
  },
};

ImportantShipmentDates.propTypes = {
  requestedPickupDate: PropTypes.string,
  scheduledPickupDate: PropTypes.string,
  plannedMoveDate: PropTypes.string,
  actualMoveDate: PropTypes.string,
  requiredDeliveryDate: PropTypes.string,
  actualPickupDate: PropTypes.string,
  requestedDeliveryDate: PropTypes.string,
  scheduledDeliveryDate: PropTypes.string,
  actualDeliveryDate: PropTypes.string,
  isPPM: PropTypes.bool,
  shipmentInfo: PropTypes.shape({
    id: PropTypes.string,
    status: ShipmentStatusesOneOf,
    shipmentType: ShipmentOptionsOneOf,
    isDiversion: PropTypes.bool,
    diversionReason: PropTypes.string,
  }),
};

export default ImportantShipmentDates;
