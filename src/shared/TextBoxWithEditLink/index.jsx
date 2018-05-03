import React from 'react';
import PropTypes from 'prop-types';

import './index.css';

const TextBoxWithEditButton = () => (
  <form>
    <textarea />
    <a href="">Edit</a>
  </form>
);

TextBoxWithEditButton.propTypes = {
  className: PropTypes.string,
  children: PropTypes.node,
  selected: PropTypes.bool,
  onClick: PropTypes.func,
};

export default TextBoxWithEditButton;
