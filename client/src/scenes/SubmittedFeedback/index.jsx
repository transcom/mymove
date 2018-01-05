// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';

import IssueCards from './IssueCards';

class SubmittedFeedback extends Component {
  constructor(props) {
    super(props);
    this.state = { issues: null };
  }
  componentDidMount() {
    this.loadIssues();
  }
  render() {
    const { issues } = this.state;
    return (
      <div className="usa-grid">
        <h1>Submitted Feedback</h1>
        <IssueCards issues={issues} />
      </div>
    );
  }
  loadIssues = () => {
    fetch('/api/v1/issues')
      .then(response => response.json())
      .then(data => this.setState({ issues: data }))
      .catch(response => console.error(response));
  };
}
export default SubmittedFeedback;
