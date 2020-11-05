import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import faCheck from '@fortawesome/free-solid-svg-icons/faCheck';
import faCaretDown from '@fortawesome/free-solid-svg-icons/faCaretDown';
import './index.css';

class ComboButton extends Component {
  container = React.createRef();
  state = {
    displayDropDown: false,
  };

  componentDidMount() {
    document.addEventListener('mousedown', this.handleClickOutside);
  }

  componentWillUnmount() {
    document.removeEventListener('mousedown', this.handleClickOutside);
  }

  handleClickOutside = (event) => {
    if (this.container.current && !this.container.current.contains(event.target)) {
      this.setState({
        displayDropDown: false,
      });
    }
  };

  handleButtonClick = () => {
    this.setState((state) => {
      return {
        displayDropDown: !state.displayDropDown,
      };
    });
  };

  render() {
    const { buttonText, disabled, children, allAreApproved } = this.props;
    return (
      <span className="container combo-button" ref={this.container}>
        <button
          className={classNames('usa-button', { 'btn__approve--green': allAreApproved })}
          disabled={disabled}
          onClick={this.handleButtonClick}
        >
          {allAreApproved && <FontAwesomeIcon className="icon" icon={faCheck} />}
          {buttonText}
          {!allAreApproved && <FontAwesomeIcon className="combo-button-icon" icon={faCaretDown} />}
        </button>
        {this.state.displayDropDown && children}
      </span>
    );
  }
}

ComboButton.propTypes = {
  buttonText: PropTypes.string,
  toolTipText: PropTypes.string,
  disabled: PropTypes.bool,
};

export default ComboButton;
