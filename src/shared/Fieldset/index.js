import React from 'react';
import { string, node, bool } from 'prop-types';
import classnames from 'classnames';

import styles from './index.module.scss';

import Hint from 'src/components/Hint/index';

const Fieldset = ({ children, legend, className, legendSrOnly, legendClassName, hintText }) => {
  const classes = classnames('usa-fieldset', className);

  const legendClasses = classnames(`usa-legend ${styles['legend-container']} ${legendClassName}`, {
    'usa-sr-only': legendSrOnly,
  });

  return (
    <fieldset data-testid="fieldset" className={classes}>
      <div className={legendClasses}>
        {legend && <legend>{legend}</legend>}
        {hintText && <Hint>{hintText}</Hint>}
      </div>
      {children}
    </fieldset>
  );
};

Fieldset.propTypes = {
  children: node,
  className: string,
  legendClassName: string,
  legendSrOnly: bool,
  legend: node,
  hintText: node,
};

Fieldset.defaultProps = {
  className: '',
  legendClassName: '',
  legendSrOnly: false,
  legend: null,
  hintText: null,
};

export default Fieldset;
