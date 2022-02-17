import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

const LeftNavSection = ({ children, sectionName, isActive, onClickHandler }) => {
  return (
    <a href={`#${sectionName}`} className={classnames({ active: isActive })} onClick={onClickHandler}>
      {children}
    </a>
  );
};

LeftNavSection.propTypes = {
  children: PropTypes.node.isRequired,
  sectionName: PropTypes.string.isRequired,
  isActive: PropTypes.bool,
  onClickHandler: PropTypes.func,
};

LeftNavSection.defaultProps = {
  isActive: false,
  onClickHandler: () => {},
};

export default LeftNavSection;
