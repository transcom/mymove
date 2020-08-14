import React from 'react';
import { string, node } from 'prop-types';

const Hint = ({ className, children, ...props }) => (
  // eslint-disable-next-line react/jsx-props-no-spreading
  <div {...props} className={`usa-hint ${className}`}>
    {children}
  </div>
);

Hint.propTypes = {
  className: string,
  children: node.isRequired,
};

Hint.defaultProps = {
  className: '',
};

export default Hint;
