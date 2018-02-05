import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import FeedbackConfirmation from 'scenes/Feedback/FeedbackConfirmation';
import FeedbackForm from 'scenes/Feedback/FeedbackForm';

import { createIssue } from './ducks';

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
    this.props.createIssue(this.state.value);
  };

  render() {
    const { value, confirmationText } = this.props;
    return (
      <div className="usa-grid">
        <h1>Report a Bug!</h1>
        <FeedbackForm
          handleChange={this.handleChange}
          handleSubmit={this.handleSubmit}
          textValue={this.state.value}
        />
        <FeedbackConfirmation confirmationText={confirmationText} />
      </div>
    );
  }
}

Feedback.propTypes = {
  createIssue: PropTypes.func.isRequired,
  // value: PropTypes.string.isRequired,
  confirmationText: PropTypes.string.isRequired,
};

function mapStateToProps(state) {
  return {
    // value: state.feedback.value,
    confirmationText: state.feedback.confirmationText,
  };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createIssue }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(Feedback);
