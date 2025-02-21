import React from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';

import { makeCalculations } from './helpers';
import styles from './ServiceItemCalculations.module.scss';

import { PaymentServiceItemParam, MTOServiceItemShape } from 'types/order';
import {
  allowedServiceItemCalculations,
  SERVICE_ITEM_CALCULATION_LABELS,
  SERVICE_ITEM_CODES,
} from 'constants/serviceItems';

const times = <FontAwesomeIcon className={styles.icon} icon="times" />;

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

    return null;
  };

  const calculations = makeCalculations(
    itemCode,
    totalAmountRequested,
    serviceItemParams,
    additionalServiceItemData,
    shipmentType,
  );

  function checkItemCode(code) {
    switch (code) {
      case (SERVICE_ITEM_CODES.FSC, SERVICE_ITEM_CODES.DOSFSC, SERVICE_ITEM_CODES.DDSFSC):
        return true;
      default:
        return false;
    }
  }

  function checkForEmptyString(input) {
    return input.length > 0 ? input : '';
  }

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
        <div>
          {calculations.map((calc, index) => {
            return (
              <div data-testid="column" className={styles.col}>
                <div data-testid="row" className={styles.row}>
                  <small data-testid="label" className={styles.descriptionTitle}>
                    {calc.label}
                  </small>
                  <small data-testid="value" className={styles.value}>
                    {calc.value === null || calc.value === '' ? null : appendSign(index, calculations.length)}
                    {calc.value}
                  </small>
                </div>
                {calc.details &&
                  calc.details.map((detail) => {
                    return (
                      <div data-testid="details" className={styles.row}>
                        <small>
                          {detail.text.includes(SERVICE_ITEM_CALCULATION_LABELS.Total) || checkItemCode(calc.itemCode)
                            ? `${detail.text.substring(0, detail.text.indexOf(':'))}:`
                            : checkForEmptyString(detail.text)}
                        </small>
                        <small>
                          {detail.text.includes(SERVICE_ITEM_CALCULATION_LABELS.Total) || checkItemCode(calc.itemCode)
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
