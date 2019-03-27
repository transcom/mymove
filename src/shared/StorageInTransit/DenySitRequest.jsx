import React, { Component } from 'react';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faBan from '@fortawesome/fontawesome-free-solid/faBan';
import './StorageInTransit.css';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import StorageInTransitOfficeDenyForm from './StorageInTransitOfficeDenyForm.jsx';

export class DenySitRequest extends Component {
  state = {
    showForm: false,
  };

  openForm = () => {
    this.setState({ showForm: true });
  };

  closeForm = () => {
    this.setState({ showForm: false });
  };

  render() {
    if (this.state.showForm)
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
    return (
      <span className="deny-sit">
        <a className="deny-sit-link" onClick={this.openForm}>
          <FontAwesomeIcon className="icon" icon={faBan} />
          Deny
        </a>
      </span>
    );
  }
}

DenySitRequest.propTypes = {};

function mapStateToProps(state) {
  return {};
}
function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators({}, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(DenySitRequest);
