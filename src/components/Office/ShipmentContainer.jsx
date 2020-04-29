import React from 'react';
import classNames from 'classnames/bind';
import * as PropTypes from 'prop-types';
import styles from './shipmentContainer.module.scss';

const cx = classNames.bind(styles);

function ShipmentContainer({ children }) {
  return <div className={`${cx('shipment-container')} container container--accent--hhg`}>{children}</div>;
}

ShipmentContainer.propTypes = {
  children: PropTypes.element,
};

export default ShipmentContainer;
