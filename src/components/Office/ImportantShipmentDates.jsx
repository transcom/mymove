import React from 'react';
import classNames from 'classnames/bind';
import * as PropTypes from 'prop-types';

import styles from './importantShipmentDates.module.scss';

const cx = classNames.bind(styles);

const ImportantShipmentDates = ({ requestedPickupDate, scheduledPickupDate }) => {
  return (
    <div className={`maxw-tablet ${cx('shipment-dates-container')}`}>
      <div>
        <h4
          className={`${cx(
            'header',
          )} font-sans-2xs text-bold text-base-dark display-inline-block margin-bottom-0 margin-top-2`}
        >
          Customer requested pick up date
        </h4>
        <h4
          className={`${cx(
            'header',
          )} font-sans-2xs text-bold text-base-dark display-inline-block margin-bottom-0 margin-top-2`}
        >
          Scheduled pick up date
        </h4>
      </div>

      <hr className="border border-base-lighter margin-bottom-105" />

      <div>
        <p className={`${cx('date')} margin-top-0 display-inline-block margin-bottom-2`}>{requestedPickupDate}</p>
        <p className={`${cx('date')} margin-top-0 display-inline-block margin-bottom-2`}>{scheduledPickupDate}</p>
      </div>
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
