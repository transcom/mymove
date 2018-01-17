import React, { Component } from 'react';

import FeedbackForm from './FeedbackForm';

class Feedback extends Component {
  constructor(props) {
    super(props);
    this.state = { value: '' };
  }
  componentDidMount() {
    document.title = 'Transcom PPP: Submit Feedback';
  }
  handleChange = e => {
    this.setState({ value: e.target.value });
  };

  handleSubmit = e => {
    e.preventDefault();
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
      </div>
    );
  }
}

export default Feedback;
