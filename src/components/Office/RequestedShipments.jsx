import React from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames/bind';
import styles from './requestedShipments.module.scss';

const cx = classNames.bind(styles);

const RequestedShipments = ({ children }) => {
  return (
    <div className={`${cx('requested-shipments')}`} data-cy="requested-shipments">
      <h4>Requested shipments</h4>
      <div className={`${cx('__content')}`}>
        {children &&
          React.Children.map(children, (child, index) => (
            // eslint-disable-next-line react/no-array-index-key
            <div key={index} className={`${cx('__item')}`}>
              {child}
            </div>
          ))}
      </div>
    </div>
  );
};

RequestedShipments.propTypes = {
  children: PropTypes.oneOfType([PropTypes.element, PropTypes.arrayOf(PropTypes.element)]),
};

export default RequestedShipments;
