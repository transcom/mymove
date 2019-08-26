import React, { Component } from 'react';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faQuestionCircle from '@fortawesome/fontawesome-free-solid/faQuestionCircle';
import './index.css';

class IconWithTooltip extends Component {
  state = {
    showTooltip: false,
  };

  toggleTooltip = () => {
    this.setState({ showTooltip: !this.state.showTooltip });
  };

  render() {
    const { showTooltip } = this.state;
    const { icon, iconClassName, toolTipText, toolTipTextClassName } = this.props;

    return (
      <div style={{ display: 'inline-block' }}>
        <FontAwesomeIcon
          aria-hidden
          className={`color_blue_link ${iconClassName}`}
          icon={icon ? icon : faQuestionCircle}
          onClick={this.toggleTooltip}
        />
        {showTooltip && (
          <div className="tooltip2">
            <div className="arrow" />
            <div className={`tooltiptext2 ${toolTipTextClassName}`}>{toolTipText}</div>
          </div>
        )}
      </div>
    );
  }
}
export default IconWithTooltip;
