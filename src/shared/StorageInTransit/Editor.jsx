import React, { Component } from 'react';
import PropTypes from 'prop-types';

import './StorageInTransit.css';

import { isValid, isSubmitting, submit, hasSubmitSucceeded } from 'redux-form';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import StorageInTransitForm, { formName as StorageInTransitFormName } from './StorageInTransitForm.jsx';

export class Editor extends Component {
  state = {
    closeOnSubmit: true,
  };

  componentDidUpdate(prevProps, prevState, snapshot) {
    if (this.props.hasSubmitSucceeded && !prevProps.hasSubmitSucceeded) {
      this.props.onClose();
    }
  }

  closeForm = () => {
    this.props.onClose();
  };

  saveAndClose = () => {
    this.setState({ closeOnSubmit: true }, () => {
      this.props.submitForm();
    });
  };

  onSubmit = values => {
    this.props.updateStorageInTransit(values);
  };

  render() {
    return (
      <div className="storage-in-transit-panel-modal">
        <div className="title">Edit SIT Request</div>
        <StorageInTransitForm onSubmit={this.onSubmit} initialValues={this.props.storageInTransit} />
        <div className="usa-grid-full">
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
              Save
            </button>
          </div>
        </div>
      </div>
    );
  }
}

Editor.propTypes = {
  updateStorageInTransit: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
  storageInTransit: PropTypes.object.isRequired,
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
export default connect(mapStateToProps, mapDispatchToProps)(Editor);
