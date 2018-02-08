import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import FeedbackConfirmation from 'scenes/Feedback/FeedbackConfirmation';
import FeedbackForm from 'scenes/Feedback/FeedbackForm';

import { createIssue, updatePendingIssueValue } from './ducks';

class Feedback extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Submit Feedback';
  }

  handleChange = e => {
    this.props.updatePendingIssueValue(e.target.value);
  };

  handleSubmit = async e => {
    e.preventDefault();
    this.props.createIssue(this.props.pendingValue);
  };

  render() {
    const { pendingValue, confirmationText } = this.props;
    return (
      <div className="usa-grid">
        <h1>Report a Bug!</h1>
        <FeedbackForm
          handleChange={this.handleChange}
          handleSubmit={this.handleSubmit}
          textValue={pendingValue}
        />
        <FeedbackConfirmation confirmationText={confirmationText} />
      </div>
    );
  }
}

Feedback.propTypes = {
  createIssue: PropTypes.func.isRequired,
  updatePendingIssueValue: PropTypes.func.isRequired,
  pendingValue: PropTypes.string.isRequired,
  confirmationText: PropTypes.string.isRequired,
};

function mapStateToProps(state) {
  return {
    pendingValue: state.feedback.pendingValue,
    confirmationText: state.feedback.confirmationText,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createIssue, updatePendingIssueValue }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Feedback);
