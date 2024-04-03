import React from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';

import { makeCalculations } from './helpers';
import styles from './ServiceItemCalculations.module.scss';

import { PaymentServiceItemParam, MTOServiceItemShape } from 'types/order';
import { allowedServiceItemCalculations, SERVICE_ITEM_CALCULATION_LABELS } from 'constants/serviceItems';

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
        <div data-testid="ServiceItemCalculations">
          {calculations.map((calc, index) => {
            return (
              <div data-testid="column" key={calc.label} className={styles.col}>
                <div data-testid="row" key={calc.value} className={styles.row}>
                  <small data-testid="label" className={styles.descriptionTitle}>
                    {calc.label}
                  </small>
                  <small data-testid="value" className={styles.value}>
                    {calc.value}
                    {appendSign(index, calculations.length)}
                  </small>
                </div>
                {calc.details &&
                  calc.details.map((detail) => {
                    return (
                      <div data-testid="details" className={styles.row}>
                        <small>
                          {detail.text.includes(SERVICE_ITEM_CALCULATION_LABELS.FSCPriceDifferenceInCents)
                            ? `${SERVICE_ITEM_CALCULATION_LABELS.FSCPriceDifferenceInCents}:`
                            : detail.text}
                        </small>
                        <small>
                          {detail.text.includes(SERVICE_ITEM_CALCULATION_LABELS.FSCPriceDifferenceInCents)
                            ? detail.text.substring(detail.text.indexOf(':') + 1)
                            : ''}
                        </small>
                      </div>
                    );
                  })}
                <hr />
              </div>
            );
          })}
        </div>
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
