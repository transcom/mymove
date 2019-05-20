import React, { Component } from 'react';
import PropTypes from 'prop-types';
import './StorageInTransit.css';
import { isValid, isSubmitting, submit, hasSubmitSucceeded } from 'redux-form';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import StorageInTransitOfficeApprovalForm, {
  formName as StorageInTransitOfficeApprovalFormName,
} from './StorageInTransitOfficeApprovalForm.jsx';
import { approveStorageInTransit } from 'shared/Entities/modules/storageInTransits';
import './StorageInTransit.css';

export class ApproveSitRequest extends Component {
  componentDidUpdate(prevProps, prevState) {
    if (this.props.hasSubmitSucceeded && !prevProps.hasSubmitSucceeded) {
      this.props.onClose();
    }
  }

  closeForm = () => {
    this.props.onClose();
  };

  approveSit = () => {
    this.props.submitForm();
  };

  onSubmit = values => {
    this.props.approveStorageInTransit(this.props.storageInTransit.shipment_id, this.props.storageInTransit.id, values);
  };

  render() {
    const { storageInTransit } = this.props;
    return (
      <div className="storage-in-transit-panel-modal">
        <div data-cy="approve-sit-request-title" className="title">
          Approve SIT Request
        </div>
        <StorageInTransitOfficeApprovalForm onSubmit={this.onSubmit} initialValues={storageInTransit} />
        <div className="usa-grid-full align-center-vertical">
          <div className="usa-width-one-half">
            <p className="cancel-link">
              <a
                data-cy="storage-in-transit-approve-cancel-link"
                className="usa-button-secondary"
                onClick={this.closeForm}
              >
                Cancel
              </a>
            </p>
          </div>
          <div className="usa-width-one-half align-right">
            <button
              data-cy="storage-in-transit-approve-button"
              className="button usa-button-primary storage-in-transit-request-form-send-request-button"
              disabled={!this.props.formEnabled}
              onClick={this.approveSit}
            >
              Approve
            </button>
          </div>
        </div>
      </div>
    );
  }
}

ApproveSitRequest.propTypes = {
  onClose: PropTypes.func.isRequired,
  storageInTransit: PropTypes.object.isRequired,
  submitForm: PropTypes.func.isRequired,
  approveStorageInTransit: PropTypes.func.isRequired,
};

function mapStateToProps(state) {
  return {
    formEnabled:
      isValid(StorageInTransitOfficeApprovalFormName)(state) &&
      !isSubmitting(StorageInTransitOfficeApprovalFormName)(state),
    hasSubmitSucceeded: hasSubmitSucceeded(StorageInTransitOfficeApprovalFormName)(state),
  };
}
function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators(
    {
      submitForm: () => submit(StorageInTransitOfficeApprovalFormName),
      approveStorageInTransit,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(ApproveSitRequest);
