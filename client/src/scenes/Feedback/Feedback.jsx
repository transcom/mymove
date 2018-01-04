import React, { Component } from 'react';
import FeedbackForm from './FeedbackForm';
import FeedbackConfirmation from './FeedbackConfirmation';

class Feedback extends Component {
  constructor(props) {
    super(props);
    this.state = { value: '' };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.confirmationText = 'sample';
  }

  handleChange = e => {
    this.setState({ value: e.target.value });
  };

  handleSubmit = e => {
    e.preventDefault();
    const config = {
      method: 'POST',
      body: JSON.stringify({ issue: 'my issue' }),
    };
    fetch('http://localhost:8080/api/v1/issues', config)
      .then(response => {
        alert(response);
        console.log(response);
        this.setState({ confirmationText: 'Response received!' });
      })
      .catch(response => {
        alert(response);
        console.log(response);
        this.setState({ confirmationText: 'Error submitting feedback' });
        // console.log(this.confirmationText); this is getting set as hoped.
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
        <FeedbackConfirmation confirmationText={this.confirmationText} />
      </div>
    );
  }
}

export default Feedback;
