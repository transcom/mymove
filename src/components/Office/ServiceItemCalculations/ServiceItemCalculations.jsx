import React from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';

import { makeCalculations } from './helpers';
import styles from './ServiceItemCalculations.module.scss';

import { PaymentServiceItemParam, MTOServiceItemShape } from 'types/order';
import { allowedServiceItemCalculations } from 'constants/serviceItems';

const times = <FontAwesomeIcon className={styles.icon} icon="times" />;
const equals = <FontAwesomeIcon className={styles.icon} icon="equals" />;

const ServiceItemCalculations = ({
  itemCode,
  totalAmountRequested,
  serviceItemParams,
  additionalServiceItemData,
  tableSize,
  shipmentType,
}) => {
  if (!allowedServiceItemCalculations.includes(itemCode) || serviceItemParams.length === 0) {
    return null;
  }

  const appendSign = (index, length) => {
    if (tableSize === 'small') {
      return null;
    }

    if (index > 0 && index !== length - 1) {
      return times;
    }

    if (index === length - 1) {
      return equals;
    }

    return null;
  };

  const calculations = makeCalculations(
    itemCode,
    totalAmountRequested,
    serviceItemParams,
    additionalServiceItemData,
    shipmentType,
  );

  return (
    <div
      data-testid="ServiceItemCalculations"
      className={classnames(styles.ServiceItemCalculations, {
        [styles.ServiceItemCalculationsSmall]: tableSize === 'small',
      })}
    >
      <h4 className={styles.title}>Calculations</h4>
      <div
        data-testid="flexGrid"
        className={classnames(styles.flexGrid, {
          [styles.flexGridSmall]: tableSize === 'small',
        })}
      >
        {calculations.map((calc, index) => {
          return (
            <div data-testid="column" key={calc.label} className={styles.col}>
              <p data-testid="value" className={styles.value}>
                {appendSign(index, calculations.length)}
                {calc.value}
              </p>
              <hr />
              <div>
                <p>
                  <small data-testid="label" className={styles.descriptionTitle}>
                    {calc.label}
                  </small>
                </p>
                <ul data-testid="details" className={styles.descriptionContent}>
                  {calc.details &&
                    calc.details.map((detail) => {
                      return (
                        <li key={detail.text}>
                          <p>
                            <small style={detail.styles}>{detail.text}</small>
                          </p>
                        </li>
                      );
                    })}
                </ul>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

ServiceItemCalculations.propTypes = {
  itemCode: PropTypes.string.isRequired,
  // in cents
  totalAmountRequested: PropTypes.number.isRequired,
  serviceItemParams: PropTypes.arrayOf(PaymentServiceItemParam),
  additionalServiceItemData: MTOServiceItemShape,
  // apply small or large styling
  tableSize: PropTypes.oneOf(['small', 'large']),
  shipmentType: PropTypes.string,
};

ServiceItemCalculations.defaultProps = {
  tableSize: 'large',
  serviceItemParams: [],
  additionalServiceItemData: {},
  shipmentType: '',
};

export default ServiceItemCalculations;
