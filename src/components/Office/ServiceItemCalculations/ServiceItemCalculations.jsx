import React from 'react';
import PropTypes from 'prop-types';

import styles from './ServiceItemCalculations.module.scss';

const ServiceItemCalculations = ({ calculations }) => {
  const appendSign = (index, length) => {
    const multiplies = <span className={styles.multiplier}>X</span>;
    const equals = <span className={styles.equal}>=</span>;

    if (index > 0 && index !== length - 1) {
      return multiplies;
    }

    if (index === length - 1) {
      return equals;
    }

    return <></>;
  };

  return (
    <div className={styles.ServiceItemCalculations}>
      <div className={styles.flexGrid}>
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
};

ServiceItemCalculations.defaultProps = {};

export default ServiceItemCalculations;
