import React from 'react';
import classnames from 'classnames';
import * as PropTypes from 'prop-types';

import styles from './importantShipmentDates.module.scss';

import DataTable from 'components/DataTable/index';
import DataTableWrapper from 'components/DataTableWrapper/index';

const ImportantShipmentDates = ({
  requestedPickupDate,
  scheduledPickupDate,
  requiredDeliveryDate,
  actualPickupDate,
  requestedDeliveryDate,
  scheduledDeliveryDate,
  actualDeliveryDate,
  isPPM,
}) => {
  const emDash = '\u2014';
  return (
    <div className={classnames('maxw-tablet', styles.shipmentDatesContainer)}>
      <DataTableWrapper className="table--data-point-group">
        {!isPPM && <DataTable columnHeaders={['Required Delivery Date']} dataRow={[requiredDeliveryDate || emDash]} />}
        <DataTable
          columnHeaders={['Requested pick up date', 'Scheduled pick up date', 'Actual pick up date']}
          dataRow={[requestedPickupDate || emDash, scheduledPickupDate || emDash, actualPickupDate || emDash]}
        />
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
  requestedDeliveryDate: '',
  scheduledDeliveryDate: '',
  actualDeliveryDate: '',
  isPPM: false,
};

ImportantShipmentDates.propTypes = {
  requestedPickupDate: PropTypes.string,
  scheduledPickupDate: PropTypes.string,
  requiredDeliveryDate: PropTypes.string,
  actualPickupDate: PropTypes.string,
  requestedDeliveryDate: PropTypes.string,
  scheduledDeliveryDate: PropTypes.string,
  actualDeliveryDate: PropTypes.string,
  isPPM: PropTypes.bool,
};

export default ImportantShipmentDates;
