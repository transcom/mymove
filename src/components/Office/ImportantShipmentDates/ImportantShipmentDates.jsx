import React from 'react';
import classnames from 'classnames';
import * as PropTypes from 'prop-types';

import styles from './importantShipmentDates.module.scss';

import DataPoint from 'components/DataPoint/index';
import DataPointGroup from 'components/DataPointGroup/index';

const ImportantShipmentDates = ({ requestedPickupDate, scheduledPickupDate }) => {
  return (
    <div className={classnames('maxw-tablet', styles.shipmentDatesContainer)}>
      <DataPointGroup>
        <DataPoint
          columnHeaders={['Customer requested pick up date', 'Scheduled pick up date']}
          dataRow={[requestedPickupDate, scheduledPickupDate]}
        />
      </DataPointGroup>
    </div>
  );
};

ImportantShipmentDates.defaultProps = {
  requestedPickupDate: '',
  scheduledPickupDate: '',
};

ImportantShipmentDates.propTypes = {
  requestedPickupDate: PropTypes.string,
  scheduledPickupDate: PropTypes.string,
};

export default ImportantShipmentDates;
