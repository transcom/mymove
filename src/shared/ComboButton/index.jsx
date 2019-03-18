import React, { Component } from 'react';
import PropTypes from 'prop-types';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faCaretDown from '@fortawesome/fontawesome-free-solid/faCaretDown';
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
    let { buttonText, disabled, children } = this.props;
    return (
      <span className="container combo-button" ref={this.container}>
        <button disabled={disabled} onClick={this.handleButtonClick}>
          {buttonText}
          <FontAwesomeIcon className="combo-button-icon" icon={faCaretDown} />
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
