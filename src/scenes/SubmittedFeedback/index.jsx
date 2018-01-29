// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import Alert from 'shared/Alert';
import IssueCards from 'scenes/SubmittedFeedback/IssueCards';

import { loadIssues } from './ducks';

class SubmittedFeedback extends Component {
  constructor(props) {
    super(props);
  }
  componentDidMount() {
    document.title = 'Transcom PPP: Submitted Feedback';
    this.props.loadIssues();
  }
  render() {
    const { issues, hasError } = this.props;
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
}

SubmittedFeedback.propTypes = {
  loadIssues: PropTypes.func.isRequired,
  issues: PropTypes.array.isRequired,
  hasError: PropTypes.bool.isRequired,
};

function mapStateToProps(state) {
  return { issues: state.issues, hasError: state.hasError };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadIssues }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(SubmittedFeedback);
