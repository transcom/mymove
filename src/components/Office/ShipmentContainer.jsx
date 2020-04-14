import React from 'react';
import classNames from 'classnames/bind';
import styles from './shipmentContainer.module.scss';

const cx = classNames.bind(styles);

function ShipmentContainer({ children }) {
  return <div className={`${cx('shipment-container')} container container--accent--blue`}>{children}</div>;
}

ShipmentContainer.propTypes = {
  children: (props, propName, componentName) => {
    // eslint-disable-next-line security/detect-object-injection
    const prop = props[propName];
    let error;

    if (React.Children.count(prop) === 0) {
      error = new Error(`\`${componentName}\` requires Children.`);
    }
    return error;
  },
};

export default ShipmentContainer;
