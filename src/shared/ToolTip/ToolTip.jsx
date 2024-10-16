import React, { useEffect, useRef, useState } from 'react';
import styles from './ToolTip.module.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

const ToolTip = ({ text, position, icon, color, closeOnLeave, title, textAreaSize, style }) => {
  // this state determines if the text is visible on mousehover/leave
  const [isVisible, setIsVisible] = useState(false);
  const tooltipRef = useRef(null);
  let textStyle; // setting initial textStyle variable

  // if the position prop is passed in, this will run checks
  // if not, it will default to top
  if (!position || position === 'top') {
    textStyle = `${styles.tooltipTextTop}`;
  } else if (position === 'left') {
    textStyle = `${styles.tooltipTextLeft}`;
  } else if (position === 'right') {
    textStyle = `${styles.tooltipTextRight}`;
  } else if (position === 'bottom') {
    textStyle = `${styles.tooltipTextBottom}`;
  }

  if (textAreaSize === 'large') {
    textStyle += ` ${styles.toolTipTextAreaLarge}`;
  }

  const determineIsVisible = () => {
    setIsVisible(!isVisible);
  };

  // this will hide the tooltips when a user clicks outside of the div
  // multiple tooltips can be open at one time, but a click will hide all of them
  const handleClickOutside = (e) => {
    if (tooltipRef.current && !tooltipRef.current.contains(e.target)) {
      setIsVisible(false);
    }
  };

  const closeOnMouseLeave = () => {
    if (closeOnLeave) {
      setIsVisible(false);
    }
  };

  useEffect(() => {
    document.addEventListener('click', handleClickOutside);
    return () => {
      document.removeEventListener('click', handleClickOutside);
    };
  }, []);

  return (
    <div
      className={styles.tooltipContainer}
      data-testid="tooltip-container"
      onMouseEnter={() => setIsVisible(true)}
      onMouseLeave={() => closeOnMouseLeave()}
      onClick={() => determineIsVisible()}
      ref={tooltipRef}
      style={style}
    >
      <FontAwesomeIcon icon={icon || 'circle-question'} color={color || 'blue'} />
      {isVisible && (
        <div className={textStyle}>
          {title && <div className={styles.popoverHeader}>{title}</div>}
          <div className={styles.popoverBody}>{text}</div>
        </div>
      )}
    </div>
  );
};

export default ToolTip;
