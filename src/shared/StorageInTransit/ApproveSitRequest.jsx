import React, { Component } from 'react';
import PropTypes from 'prop-types';
import './StorageInTransit.css';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import StorageInTransitOfficeApprovalForm from './StorageInTransitOfficeApprovalForm.jsx';
import './StorageInTransit.css';

export class ApproveSitRequest extends Component {
  closeForm = () => {
    this.props.onClose();
  };

  render() {
    const { storageInTransit } = this.props;
    return (
      <div className="storage-in-transit-panel-modal">
        <div className="title">Approve SIT Request</div>
        <StorageInTransitOfficeApprovalForm initialValues={storageInTransit} />
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
};

function mapStateToProps(state) {
  return {};
}
function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators({}, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(ApproveSitRequest);
