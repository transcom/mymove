import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import Alert from 'shared/Alert';
import { reduxifyForm } from 'shared/JsonSchemaForm';

import { createIssue } from './ducks';

const FeedbackForm = reduxifyForm('Feedback');

export class Feedback extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Submit Feedback';
  }

  handleSubmit = values => {
    this.props.createIssue(values);
  };

  render() {
    const { hasErrored, hasSucceeded } = this.props;
    return (
      <div className="usa-grid">
        <h1>Report a Bug!</h1>
        <FeedbackForm
          onSubmit={this.handleSubmit}
          schema={this.props.schema}
          uiSchema={this.props.uiSchema}
        />
        {hasErrored && (
          <Alert type="error" heading="Submission Error">
            Something went wrong
          </Alert>
        )}
        {hasSucceeded && (
          <Alert type="success" heading="Submission Successful">
            Your issue is submitted.
          </Alert>
        )}
      </div>
    );
  }
}

Feedback.propTypes = {
  createIssue: PropTypes.func.isRequired,
  schema: PropTypes.object.isRequired,
  uiSchema: PropTypes.object.isRequired,
  hasErrored: PropTypes.bool.isRequired,
  hasSucceeded: PropTypes.bool.isRequired,
};

function mapStateToProps(state) {
  const props = { ...state.feedback, schema: {} };
  if (state.swagger.spec) {
    props.schema = state.swagger.spec.definitions.CreateIssuePayload;
  }
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createIssue }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Feedback);
