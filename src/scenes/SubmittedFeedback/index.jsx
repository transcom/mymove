// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';

import IssueCards from 'scenes/SubmittedFeedback/IssueCards';
import Alert from 'shared/Alert';
import { IssuesIndex } from 'shared/api.js';

class SubmittedFeedback extends Component {
  constructor(props) {
    super(props);
    this.state = { issues: null, hasError: false };
  }
  componentDidMount() {
    document.title = 'Transcom PPP: Submitted Feedback';
    this.loadIssues();
  }
  render() {
    const { issues, hasError } = this.state;
    return (
      <div className="usa-grid">
        <h1>Submitted Feedback</h1>
        {hasError && (
          <Alert type="error" heading="Server Error">
            There was a problem loading the issues from the server.
          </Alert>
        )}
        {!hasError && <IssueCards issues={issues} />}
      </div>
    );
  }
  loadIssues = async () => {
    try {
      const issues = await IssuesIndex();
      this.setState({ issues });
    } catch (e) {
      //componentDidCatch will not get fired because this is async
      //todo: how to we want to monitor errors
      console.error(e);
      this.setState({ hasError: true });
    }
  };
}
export default SubmittedFeedback;
