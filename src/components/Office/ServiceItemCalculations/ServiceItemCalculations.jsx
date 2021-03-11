import React from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './ServiceItemCalculations.module.scss';

const ServiceItemCalculations = ({ calculations, tableSize }) => {
  const appendSign = (index, length) => {
    const multiplies = <FontAwesomeIcon className={styles.icon} icon="times" />;
    const equals = <FontAwesomeIcon className={styles.icon} icon="equals" />;

    if (index > 0 && index !== length - 1) {
      return multiplies;
    }

    if (index === length - 1) {
      return equals;
    }

    return <></>;
  };

  return (
    <div
      className={`${styles.ServiceItemCalculations}
        ${tableSize === 'small' ? styles.ServiceItemCalculationsSmall : ''}`}
    >
      <h4 className={styles.title}>Calculations</h4>
      <div className={`${styles.flexGrid} ${tableSize === 'small' ? styles.flexGridSmall : ''}`}>
        {calculations.map((calc, index) => {
          return (
            <div key={calc.label} className={styles.col}>
              <div className={styles.value}>
                {appendSign(index, calculations.length)}
                {calc.value}
              </div>
              <hr />
              <div>
                <div className={styles.descriptionTitle}>{calc.label}</div>
                <div className={styles.descriptionContent}>
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
