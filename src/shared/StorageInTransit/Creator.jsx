import React, { Component } from 'react';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';

import './StorageInTransit.css';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import StorageInTransitForm from './StorageInTransitForm';

export class Creator extends Component {
  state = { showForm: false };

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
          <div className="title">Request SIT</div>
          <StorageInTransitForm />
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
              >
                Send Request
              </button>
            </div>
          </div>
        </div>
      );
    return (
      <div className="add-request storage-in-transit-hr-top">
        <a onClick={this.openForm}>
          <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} />
          Request SIT
        </a>
      </div>
    );
  }
}

Creator.propTypes = {};

function mapStateToProps(state) {
  return {};
}

function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators({}, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(Creator);
