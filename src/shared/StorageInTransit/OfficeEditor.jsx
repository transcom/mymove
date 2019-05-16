import React, { Component } from 'react';
import PropTypes from 'prop-types';

import './StorageInTransit.css';

import { isValid, isSubmitting, submit, hasSubmitSucceeded } from 'redux-form';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { approveStorageInTransit, denyStorageInTransit } from 'shared/Entities/modules/storageInTransits';
import StorageInTransitOfficeEditForm, {
  formName as StorageInTransitOfficeEditFormName,
} from './StorageInTransitOfficeEditForm.jsx';

export class OfficeEditor extends Component {
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
    if (this.props.storageInTransit.status === 'APPROVED') {
      this.props.approveStorageInTransit(
        this.props.storageInTransit.shipment_id,
        this.props.storageInTransit.id,
        values,
      );
    }
    if (this.props.storageInTransit.status === 'DENIED') {
      this.props.denyStorageInTransit(this.props.storageInTransit.shipment_id, this.props.storageInTransit.id, values);
    }
  };

  render() {
    return (
      <div className="storage-in-transit-panel-modal">
        <div className="editable-panel is-editable">
          <div data-cy="sit-authorization-title" className="sit-authorization title">
            Edit SIT authorization
          </div>
          <StorageInTransitOfficeEditForm onSubmit={this.onSubmit} initialValues={this.props.storageInTransit} />
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
                data-cy="sit-editor-save-button"
                className="button usa-button-primary"
                disabled={!this.props.formEnabled}
                onClick={this.saveAndClose}
              >
                Save
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

OfficeEditor.propTypes = {
  approveStorageInTransit: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
  storageInTransit: PropTypes.object.isRequired,
  submitForm: PropTypes.func.isRequired,
  denyStorageInTransit: PropTypes.func.isRequired,
};

function mapStateToProps(state) {
  return {
    formEnabled:
      isValid(StorageInTransitOfficeEditFormName)(state) && !isSubmitting(StorageInTransitOfficeEditFormName)(state),
    hasSubmitSucceeded: hasSubmitSucceeded(StorageInTransitOfficeEditFormName)(state),
  };
}

function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators(
    {
      submitForm: () => submit(StorageInTransitOfficeEditFormName),
      approveStorageInTransit,
      denyStorageInTransit,
    },
    dispatch,
  );
}
export default connect(mapStateToProps, mapDispatchToProps)(OfficeEditor);
