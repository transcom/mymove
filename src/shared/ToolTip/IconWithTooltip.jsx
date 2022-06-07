import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import styles from './IconWithTooltip.module.scss';

class IconWithTooltip extends Component {
  state = {
    showTooltip: false,
  };

  toggleTooltip = () => {
    this.setState((prevState) => ({ showTooltip: !prevState.showTooltip }));
  };

  render() {
    const { showTooltip } = this.state;
    const { icon, iconClassName, toolTipText, toolTipTextClassName, toolTipStyles } = this.props;

    return (
      <div style={{ display: 'inline-block' }}>
        <FontAwesomeIcon
          aria-hidden
          className={`${styles['color_blue_link']} ${iconClassName}`}
          icon={icon ? icon : 'circle-question'}
          onClick={this.toggleTooltip}
        />
        {showTooltip && (
          <div data-testid="tooltip" className={styles['tooltip']} style={{ ...toolTipStyles }}>
            <div className={styles['arrow']} onClick={this.toggleTooltip} />
            <div className={`${styles['tooltiptext']} ${toolTipTextClassName}`}>{toolTipText}</div>
          </div>
        )}
      </div>
    );
  }
}

IconWithTooltip.propTypes = {
  icon: PropTypes.node,
  iconClassName: PropTypes.string,
  toolTipText: PropTypes.string.isRequired,
  toolTipTextClassName: PropTypes.string,
  toolTipStyles: PropTypes.object,
};

export default IconWithTooltip;
