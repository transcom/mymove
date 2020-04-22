import React from 'react';
import classNames from 'classnames/bind';
import * as PropTypes from 'prop-types';
import styles from './importantShipmentDates.module.scss';

const cx = classNames.bind(styles);

const ImportantShipmentDates = ({ requestedPickupDate, scheduledPickupDate }) => {
  return (
    <div className={`container container--gray maxw-tablet ${cx('shipment-dates-container')}`}>
      <div>
        <p className={`${cx('header')} font-sans-2xs text-bold text-base-dark display-inline-block margin-bottom-0`}>
          Customer requested pick up date
        </p>
        <p className={`${cx('header')} font-sans-2xs text-bold text-base-dark display-inline-block margin-bottom-0`}>
          Scheduled pick up date
        </p>
      </div>

      <hr className="border border-base-lighter" />

      <div>
        <p className={`${cx('date')} margin-top-0 display-inline-block`}>{requestedPickupDate}</p>
        <p className={`${cx('date')} margin-top-0 display-inline-block`}>{scheduledPickupDate}</p>
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
