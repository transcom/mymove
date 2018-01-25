import React, { Component } from 'react';

import FeedbackForm from 'scenes/Feedback/FeedbackForm';
import FeedbackConfirmation from 'scenes/Feedback/FeedbackConfirmation';
import { CreateIssue } from 'shared/api.js';
class Feedback extends Component {
  constructor(props) {
    super(props);
    this.state = {
      value: '',
      confirmationText: '',
    };
  }
  componentDidMount() {
    document.title = 'Transcom PPP: Submit Feedback';
  }
  handleChange = e => {
    this.setState({ value: e.target.value });
  };

  handleSubmit = async e => {
    e.preventDefault();

    try {
      await CreateIssue(this.state.value);
      this.setState({ confirmationText: 'Response received!' });
    } catch (e) {
      //todo: how do we want to monitor errors
      this.setState({ confirmationText: 'Error submitting feedback' });
    }
  };

  render() {
    return (
      <div className="usa-grid">
        <h1>Report a Bug!</h1>
        <FeedbackForm
          handleChange={this.handleChange}
          handleSubmit={this.handleSubmit}
          textValue={this.state.value}
        />
        <FeedbackConfirmation confirmationText={this.state.confirmationText} />
      </div>
    );
  }
}

export default Feedback;
