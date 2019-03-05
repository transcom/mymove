import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';

// TODO uncomment. couldn't run tests with this import for some reason
// import 'shared/shared.css'
import './index.css';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faCaretDown from '@fortawesome/fontawesome-free-solid/faCaretDown';

function styleListItems(items) {
  const isDisabled = item => (item.disabled ? 'disabled' : '');
  const isLastItem = (item, index) => (index === items.length - 1 ? 'last-item' : '');
  const liClasses = (item, index) => classNames(isDisabled(item), isLastItem(item, index));
  return items.map((item, index) => <li className={liClasses(item, index)}>{item.value}</li>);
}

class DropDown extends Component {
  render() {
    let { items } = this.props;
    return (
      <div className="dropdown">
        <ul className="dropdown">{styleListItems(items)}</ul>
      </div>
    );
  }
}

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
      <span className="container combo-button tooltip" ref={this.container}>
        <button disabled={disabled} onClick={this.handleButtonClick}>
          {buttonText}
          <FontAwesomeIcon className="combo-button-icon" icon={faCaretDown} />
        </button>
        {toolTipText && disabled && <span className="tooltiptext tooltiptext-large">{toolTipText}</span>}
        {this.state.displayDropDown && <DropDown items={items} />}
      </span>
    );
  }
}

DropDown.propTypes = {
  items: PropTypes.array,
};

ComboButton.propTypes = {
  buttonText: PropTypes.string,
  toolTipText: PropTypes.string,
  disabled: PropTypes.bool,
};

export default ComboButton;
