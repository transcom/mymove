import React, { Component } from 'react';
import FeedbackForm from './FeedbackForm';
import FeedbackConfirmation from './FeedbackConfirmation';

class Feedback extends Component {
  constructor(props) {
    super(props);
    this.state = {
      value: '',
      confirmationText: '',
    };
  }

  handleChange = e => {
    this.setState({ value: e.target.value });
  };

  handleSubmit = e => {
    e.preventDefault();
    const config = {
      method: 'POST',
      headers: {
        'Content-Type': 'text/plain',
      },
      body: JSON.stringify({ body: this.state.value }),
    };
    fetch('/api/v1/issues', config).then(response => {
      if (response.ok) {
        this.setState({ confirmationText: 'Response received!' });
      } else {
        this.setState({ confirmationText: 'Error submitting feedback' });
      }
    });
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
