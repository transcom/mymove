import React, { Component } from 'react';
import PropTypes from 'prop-types';
import './StorageInTransit.css';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import StorageInTransitOfficeDenyForm from './StorageInTransitOfficeDenyForm.jsx';

export class DenySitRequest extends Component {
  closeForm = () => {
    this.props.onClose();
  };

  render() {
    return (
      <div className="storage-in-transit-panel-modal">
        <div className="title">Deny SIT Request</div>
        <StorageInTransitOfficeDenyForm />
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
};

function mapStateToProps(state) {
  return {};
}
function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators({}, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(DenySitRequest);
