import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import FeedbackConfirmation from 'scenes/Feedback/FeedbackConfirmation';
import FeedbackForm from 'scenes/Feedback/FeedbackForm';

import { createIssue } from './ducks';

class Feedback extends Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
    document.title = 'Transcom PPP: Submit Feedback';
    this.props.createIssue; // I think this shouldn't go here because it's used to create something, not display something that already exists.
  }
  handleChange = e => {
    this.setState({ value: e.target.value }); // Hmmm
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
          textValue={this.state.value}
        />
        <FeedbackConfirmation confirmationText={this.props.confirmationText} />
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
    value: state.Feedback.value, // These two are guesses/placeholders
    confirmationText: state.Feedback.confirmationText, // and feel wonky
  };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createIssue }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Feedback);
