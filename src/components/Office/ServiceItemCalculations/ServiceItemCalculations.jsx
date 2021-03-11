import React from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './ServiceItemCalculations.module.scss';

const ServiceItemCalculations = ({ calculations, tableSize }) => {
  const appendSign = (index, length) => {
    if (tableSize === 'small') {
      return <></>;
    }

    const times = <FontAwesomeIcon className={styles.icon} icon="times" />;
    const equals = <FontAwesomeIcon className={styles.icon} icon="equals" />;

    if (index > 0 && index !== length - 1) {
      return times;
    }

    if (index === length - 1) {
      return equals;
    }

    return <></>;
  };

  return (
    <div
      data-testid="ServiceItemCalculations"
      className={`${styles.ServiceItemCalculations}
        ${tableSize === 'small' ? styles.ServiceItemCalculationsSmall : ''}`}
    >
      <h4 className={styles.title}>Calculations</h4>
      <div data-testid="flexGrid" className={`${styles.flexGrid} ${tableSize === 'small' ? styles.flexGridSmall : ''}`}>
        {calculations.map((calc, index) => {
          return (
            <div data-testid="column" key={calc.label} className={styles.col}>
              <div data-testid="value" className={styles.value}>
                {appendSign(index, calculations.length)}
                {calc.value}
              </div>
              <hr />
              <div>
                <div data-testid="label" className={styles.descriptionTitle}>
                  {calc.label}
                </div>
                <div data-testid="details" className={styles.descriptionContent}>
                  {calc.details &&
                    calc.details.map((detail, i) => {
                      if (i === calc.details.length - 1) {
                        return <React.Fragment key={detail}>{detail}</React.Fragment>;
                      }

                      // each item, add line breaks
                      return (
                        <React.Fragment key={detail}>
                          {detail} <br />
                        </React.Fragment>
                      );
                    })}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

ServiceItemCalculations.propTypes = {
  // collection of ordered calculations and last item is the Total amount requested
  calculations: PropTypes.arrayOf(
    PropTypes.shape({
      value: PropTypes.string.isRequired,
      label: PropTypes.string.isRequired,
      details: PropTypes.arrayOf(PropTypes.string),
    }),
  ).isRequired,
  // apply small or large styling
  tableSize: PropTypes.oneOf(['small', 'large']),
};

ServiceItemCalculations.defaultProps = {
  tableSize: 'large',
};

export default ServiceItemCalculations;
