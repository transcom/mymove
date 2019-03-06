import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faCaretDown from '@fortawesome/fontawesome-free-solid/faCaretDown';
import ToolTip from 'shared/ToolTip';
import './index.css';

class ComboButton extends Component {
  container = React.createRef();
  state = {
    displayDropDown: false,
  };

  handleClickOutside = event => {
    if (this.container.current && !this.container.current.contains(event.target)) {
      this.setState({
        displayDropDown: false,
      });
    }
  };

  componentDidMount() {
    document.addEventListener('mousedown', this.handleClickOutside);
  }
  componentWillUnmount() {
    document.removeEventListener('mousedown', this.handleClickOutside);
  }

  handleButtonClick = () => {
    this.setState(state => {
      return {
        displayDropDown: !state.displayDropDown,
      };
    });
  };

  render() {
    let { buttonText, toolTipText, disabled, items = [] } = this.props;
    return (
      <span className="container combo-button" ref={this.container}>
        <ToolTip text={toolTipText} disabled={!disabled} textStyle={'tooltiptext-large'}>
          <button disabled={disabled} onClick={this.handleButtonClick}>
            {buttonText}
            <FontAwesomeIcon className="combo-button-icon" icon={faCaretDown} />
          </button>
          {this.state.displayDropDown && <DropDown items={items} />}
        </ToolTip>
      </span>
    );
  }
}

const styleListItems = items => {
  const liClasses = item => classNames({ disabled: item.disabled });
  return items.map(item => <li className={liClasses(item)}>{item.value}</li>);
};

const DropDown = props => {
  let { items } = props;
  return (
    <div className="dropdown">
      <ul className="dropdown">{styleListItems(items)}</ul>
    </div>
  );
};

DropDown.propTypes = {
  items: PropTypes.array,
};

ComboButton.propTypes = {
  buttonText: PropTypes.string,
  toolTipText: PropTypes.string,
  disabled: PropTypes.bool,
};

export default ComboButton;
