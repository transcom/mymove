/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import PropTypes from 'prop-types';

const minWidth = ({ children }) => (
  <div style="min-width:1240px; overflow-x:scroll">
    {children}
  </div>
);

minWidth.propTypes = {
  children: PropTypes.node.isRequired,
};

export default minWidth;
