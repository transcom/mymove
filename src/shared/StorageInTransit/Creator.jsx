import React, { Component } from 'react';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';
import PropTypes from 'prop-types';

import './StorageInTransit.css';

import { isValid, isSubmitting, submit, hasSubmitSucceeded } from 'redux-form';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import StorageInTransitForm, { formName as StorageInTransitFormName } from './StorageInTransitForm.jsx';

export class Creator extends Component {
  state = {
    showForm: false,
    closeOnSubmit: true,
  };

  componentDidUpdate(prevProps, prevState, snapshot) {
    if (this.props.hasSubmitSucceeded && !prevProps.hasSubmitSucceeded) {
      this.setState({ showForm: false });
    }
  }

  openForm = () => {
    this.setState({ showForm: true });
  };

  closeForm = () => {
    this.setState({ showForm: false });
  };

  saveAndClose = () => {
    this.setState({ closeOnSubmit: true }, () => {
      this.props.submitForm();
      this.props.onFormActivation(false);
    });
  };

  onSubmit = values => {
    this.props.saveStorageInTransit(values);
  };

  render() {
    if (this.state.showForm)
      return (
        <div className="storage-in-transit-panel-modal">
          <div className="title">Request SIT</div>
          <StorageInTransitForm onSubmit={this.onSubmit} />
          <div className="usa-grid-full align-center-vertical">
            <div className="usa-width-one-half">
              <p className="cancel-link">
                <a className="usa-button-secondary" onClick={this.closeForm}>
                  Cancel
                </a>
              </p>
            </div>
            <div className="usa-width-one-half align-right">
              <button
                className="button usa-button-primary storage-in-transit-request-form-send-request-button"
                disabled={!this.props.formEnabled}
                onClick={this.saveAndClose}
              >
                Send Request
              </button>
            </div>
          </div>
        </div>
      );
    return (
      <div className="add-request">
        <a onClick={this.openForm}>
          <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} />
          Request SIT
        </a>
      </div>
    );
  }
}

Creator.propTypes = {
  formEnabled: PropTypes.bool.isRequired,
  onFormActivation: PropTypes.func.isRequired,
  saveStorageInTransit: PropTypes.func.isRequired,
  submitForm: PropTypes.func.isRequired,
};

function mapStateToProps(state) {
  return {
    formEnabled: isValid(StorageInTransitFormName)(state) && !isSubmitting(StorageInTransitFormName)(state),
    hasSubmitSucceeded: hasSubmitSucceeded(StorageInTransitFormName)(state),
  };
}

function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators(
    {
      submitForm: () => submit(StorageInTransitFormName),
    },
    dispatch,
  );
}
export default connect(mapStateToProps, mapDispatchToProps)(Creator);
