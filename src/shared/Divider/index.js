import React from 'react';
import { string } from 'prop-types';

const Divider = ({ className }) => <hr className={className} />;

Divider.propTypes = {
  className: string,
};

Divider.defaultProps = {
  className: '',
};

export default Divider;
