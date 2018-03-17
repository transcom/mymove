import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import FeedbackConfirmation from 'scenes/Feedback/FeedbackConfirmation';
import { reduxifyForm } from 'shared/JsonSchemaForm';

import { createIssue } from './ducks';
import { loadSchema } from 'shared/Swagger/ducks';

const FeedbackForm = reduxifyForm('Feedback');

export class Feedback extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Submit Feedback';
    this.props.loadSchema();
  }

  handleSubmit = values => {
    this.props.createIssue(values);
  };

  render() {
    const { confirmationText } = this.props;
    return (
      <div className="usa-grid">
        <h1>Report a Bug!</h1>
        <FeedbackForm
          onSubmit={this.handleSubmit}
          schema={this.props.schema}
          uiSchema={this.props.uiSchema}
        />
        <FeedbackConfirmation confirmationText={confirmationText} />
      </div>
    );
  }
}

Feedback.propTypes = {
  createIssue: PropTypes.func.isRequired,
  confirmationText: PropTypes.string.isRequired,
  loadSchema: PropTypes.func.isRequired,
  schema: PropTypes.object.isRequired,
  uiSchema: PropTypes.object.isRequired,
  hasSchemaError: PropTypes.bool.isRequired,
  hasSubmitError: PropTypes.bool.isRequired,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapStateToProps(state) {
  const props = { ...state.feedback, schema: {} };
  if (state.swagger.spec) {
    props.schema = state.swagger.spec.definitions.CreateIssuePayload;
  }
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadSchema, createIssue }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Feedback);
