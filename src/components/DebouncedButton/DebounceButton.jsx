import React, { useRef } from 'react';
import { Button } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';

const DebounceButton = ({ onClick, delay, ariaLabel, chidren, ...props }) => {
  const lastClick = useRef(0);

  const handleClick = (e) => {
    const now = Date.now();
    if (now - lastClick.current > delay) {
      lastClick.current = now;
      onClick?.(e);
    }
  };

  // eslint-disable-next-line react/jsx-props-no-spreading
  return <Button datatest-id="debounce-button" onClick={handleClick} {...props} />;
};

DebounceButton.defaultProps = {
  delay: 2000,
};

DebounceButton.propTypes = {
  children: PropTypes.node.isRequired,
  onClick: PropTypes.func.isRequired,
  delay: PropTypes.number,
};

export default DebounceButton;
