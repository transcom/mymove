import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import FeedbackConfirmation from 'scenes/Feedback/FeedbackConfirmation';
import FeedbackForm from 'scenes/Feedback/FeedbackForm';

import { createIssue } from './ducks';

class Feedback extends Component {
  // constructor(props) {
  //   super(props),
  //   // this.props.
  //   value = '',
  // }

  componentDidMount() {
    document.title = 'Transcom PPP: Submit Feedback';
  }

  handleChange = e => {
    value = e.target.value; // Value should be gotten in this way, but it should be assigned to something else in another way. Assigned to a variable that's referenced in mapStateToProps?
  };

  handleSubmit = async e => {
    e.preventDefault();
    createIssue(); // Is this enough?
  };

  render() {
    const { value, confirmationText } = this.props;
    return (
      <div className="usa-grid">
        <h1>Report a Bug!</h1>
        <FeedbackForm
          handleChange={this.handleChange}
          handleSubmit={this.handleSubmit}
          textValue={value}
        />
        <FeedbackConfirmation confirmationText={confirmationText} />
      </div>
    );
  }
}

Feedback.propTypes = {
  createIssue: PropTypes.func.isRequired,
  value: PropTypes.string.isRequired,
  confirmationText: PropTypes.string.isRequired,
};

function mapStateToProps(state) {
  return {
    value: value, // These two are guesses/placeholders
    confirmationText: state.Feedback.confirmationText, // and feel wonky
  };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createIssue }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Feedback);
