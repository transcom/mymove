import React, { useState } from 'react';
import styles from './HoverToolTip.module.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

const HoverTooltip = ({ text, position, icon }) => {
  const [isVisible, setIsVisible] = useState(false);
  let textStyle;

  if (!position || position === 'top') {
    textStyle = `${styles.tooltipTextTop}`;
  } else if (position === 'left') {
    textStyle = `${styles.tooltipTextLeft}`;
  } else if (position === 'right') {
    textStyle = `${styles.tooltipTextRight}`;
  } else if (position === 'bottom') {
    textStyle = `${styles.tooltipTextBottom}`;
  }

  return (
    <div
      className={styles.tooltipContainer}
      onMouseEnter={() => setIsVisible(true)}
      onMouseLeave={() => setIsVisible(false)}
    >
      <FontAwesomeIcon icon={icon || 'circle-question'} />
      {isVisible && <div className={textStyle}>{text}</div>}
    </div>
  );
};

export default HoverTooltip;
