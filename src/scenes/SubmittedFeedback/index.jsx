// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';

import IssueCards from 'scenes/SubmittedFeedback/IssueCards';
import { IssuesIndex } from 'shared/api.js';

class SubmittedFeedback extends Component {
  constructor(props) {
    super(props);
    this.state = { issues: null };
  }
  componentDidMount() {
    this.loadIssues();
    document.title = 'Transcom PPP: Submitted Feedback';
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
    IssuesIndex().then(response => {
      if (response.ok) {
        this.setState({ issues: response.body });
      } else {
        console.error(response);
      }
    });
  };
}
export default SubmittedFeedback;
