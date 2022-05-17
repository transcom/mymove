import React from 'react';
import classnames from 'classnames';
import * as PropTypes from 'prop-types';

import styles from './importantShipmentDates.module.scss';

import DataTable from 'components/DataTable/index';
import DataTableWrapper from 'components/DataTableWrapper/index';

const ImportantShipmentDates = ({ requestedPickupDate, scheduledPickupDate, requiredDeliveryDate }) => {
  const emDash = '\u2014';
  return (
    <div className={classnames('maxw-tablet', styles.shipmentDatesContainer)}>
      <DataTableWrapper className="table--data-point-group">
        <DataTable columnHeaders={['Customer requested pick up date']} dataRow={[requestedPickupDate || emDash]} />
        <DataTable
          columnHeaders={['Scheduled pick up date', 'Required delivery date']}
          dataRow={[scheduledPickupDate || emDash, requiredDeliveryDate || emDash]}
        />
      </DataTableWrapper>
    </div>
  );
};

ImportantShipmentDates.defaultProps = {
  requestedPickupDate: '',
  scheduledPickupDate: '',
  requiredDeliveryDate: '',
};

ImportantShipmentDates.propTypes = {
  requestedPickupDate: PropTypes.string,
  scheduledPickupDate: PropTypes.string,
  requiredDeliveryDate: PropTypes.string,
};

export default ImportantShipmentDates;
