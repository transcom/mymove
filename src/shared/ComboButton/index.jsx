import React, { Component } from 'react';
import PropTypes from 'prop-types';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import faCaretDown from '@fortawesome/fontawesome-free-solid/faCaretDown';
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

  handleClickOutside = event => {
    if (this.container.current && !this.container.current.contains(event.target)) {
      this.setState({
        displayDropDown: false,
      });
    }
  };

  handleButtonClick = () => {
    this.setState(state => {
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
          className={allAreApproved ? 'btn__approve--green' : ''}
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
