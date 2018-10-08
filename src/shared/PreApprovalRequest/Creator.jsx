import React, { Component } from 'react';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';
import PropTypes from 'prop-types';
import PreApprovalRequestForm, {
  formName as PreApprovalRequestFormName,
} from 'shared/PreApprovalRequest/PreApprovalForm.jsx';
import {
  submit,
  isValid,
  isSubmitting,
  reset,
  hasSubmitSucceeded,
} from 'redux-form';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
export class Creator extends Component {
  state = { showForm: false, closeOnSubmit: true };
  componentDidUpdate(prevProps, prevState, snapshot) {
    if (this.props.hasSubmitSucceeded && !prevProps.hasSubmitSucceeded)
      if (this.state.closeOnSubmit) this.setState({ showForm: false });
      else this.props.clearForm();
  }
  openForm = () => {
    this.setState({ showForm: true });
  };
  closeForm = () => {
    this.setState({ showForm: false });
  };
  onSubmit = values => {
    this.props.savePreApprovalRequest(values);
  };
  saveAndClear = () => {
    this.setState({ closeOnSubmit: false }, () => {
      this.props.submitForm();
    });
  };
  saveAndClose = () => {
    this.setState({ closeOnSubmit: true }, () => {
      this.props.submitForm();
    });
  };
  render() {
    if (this.state.showForm)
      return (
        <div className="accessorial-panel-modal">
          <div className="title">Add a request</div>
          <PreApprovalRequestForm
            accessorials={this.props.accessorials}
            onSubmit={this.onSubmit}
          />
          <div className="usa-grid">
            <div className="usa-width-one-half cancel">
              <button className="usa-button-secondary" onClick={this.closeForm}>
                Cancel
              </button>
            </div>

            <div className="usa-width-one-half">
              <button
                className="button button-secondary"
                disabled={!this.props.formEnabled}
                onClick={this.saveAndClear}
              >
                Save &amp; Add Another
              </button>
              <button
                className="button button-primary"
                disabled={!this.props.formEnabled}
                onClick={this.saveAndClose}
              >
                Save &amp; Close
              </button>
            </div>
          </div>
        </div>
      );
    return (
      <div className="add-request">
        <a onClick={this.openForm}>
          <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} />
          Add a request
        </a>
      </div>
    );
  }
}
Creator.propTypes = {
  accessorials: PropTypes.array,
  savePreApprovalRequest: PropTypes.func.isRequired,
  formEnabled: PropTypes.bool.isRequired,
  hasSubmitSucceeded: PropTypes.bool.isRequired,
  submitForm: PropTypes.func.isRequired,
  clearForm: PropTypes.func.isRequired,
};

function mapStateToProps(state) {
  return {
    formEnabled:
      isValid(PreApprovalRequestFormName)(state) &&
      !isSubmitting(PreApprovalRequestFormName)(state),
    hasSubmitSucceeded: hasSubmitSucceeded(PreApprovalRequestFormName)(state),
  };
}

function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators(
    {
      submitForm: () => submit(PreApprovalRequestFormName),
      clearForm: () => reset(PreApprovalRequestFormName),
    },
    dispatch,
  );
}
export default connect(mapStateToProps, mapDispatchToProps)(Creator);
