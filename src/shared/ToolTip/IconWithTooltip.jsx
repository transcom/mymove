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
    const { icon, iconClassName, toolTipText, toolTipClassName } = this.props;

    return (
      <div style={{ display: 'inline-block' }}>
        <FontAwesomeIcon
          aria-hidden
          className={`color_blue_link ${iconClassName}`}
          icon={icon ? icon : faQuestionCircle}
          onClick={this.toggleTooltip}
        />
        <div className="tooltip">
          <span className="tooltiptext">{toolTipText}</span>
        </div>
      </div>
    );
  }
}
export default IconWithTooltip;
