import React, { Component } from 'react';
import { isValid, isSubmitting, submit, hasSubmitSucceeded } from 'redux-form';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import { denyStorageInTransit } from 'shared/Entities/modules/storageInTransits';

import './StorageInTransit.css';

import StorageInTransitOfficeDenyForm, {
  formName as StorageInTransitOfficeDenyFormName,
} from './StorageInTransitOfficeDenyForm.jsx';

export class DenySitRequest extends Component {
  closeForm = () => {
    this.props.onClose();
  };

  denySit = () => {
    this.props.submitForm();
  };

  onSubmit = values => {
    this.props.denyStorageInTransit(this.props.storageInTransit.shipment_id, this.props.storageInTransit.id, values);
  };

  componentDidUpdate(prevProps) {
    if (this.props.hasSubmitSucceeded && !prevProps.hasSubmitSucceeded) {
      this.props.onClose();
    }
  }

  render() {
    return (
      <div className="storage-in-transit-panel-modal">
        <div className="title">Deny SIT Request</div>
        <StorageInTransitOfficeDenyForm onSubmit={this.onSubmit} initialValues={this.props.storageInTransit} />
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
              className="button usa-button-primary"
              data-cy="storage-in-transit-deny-button"
              disabled={!this.props.formEnabled}
              onClick={this.denySit}
            >
              Deny
            </button>
          </div>
        </div>
      </div>
    );
  }
}

DenySitRequest.propTypes = {
  onClose: PropTypes.func.isRequired,
  storageInTransit: PropTypes.object.isRequired,
  submitForm: PropTypes.func.isRequired,
  formEnabled: PropTypes.bool,
  hasSubmitSucceeded: PropTypes.bool,
  denyStorageInTransit: PropTypes.func.isRequired,
};

function mapStateToProps(state) {
  return {
    formEnabled:
      isValid(StorageInTransitOfficeDenyFormName)(state) && !isSubmitting(StorageInTransitOfficeDenyFormName)(state),
    hasSubmitSucceeded: hasSubmitSucceeded(StorageInTransitOfficeDenyFormName)(state),
  };
}
function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators(
    {
      submitForm: () => submit(StorageInTransitOfficeDenyFormName),
      denyStorageInTransit,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(DenySitRequest);
